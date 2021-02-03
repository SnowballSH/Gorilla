package vm

import (
	"Gorilla/code"
	"Gorilla/object"
	"fmt"
)

const StackSize = 1 << 12

type VM struct {
	Constants []object.BaseObject
	Messages  []object.Message
	mp        int

	Instructions []code.Opcode
	ip           int

	Stack []object.BaseObject
	sp    int // Always points to the next value. Top of stack is stack[sp-1]

	LastPopped object.BaseObject

	Env *object.Environment
}

func New(bytecodes []code.Opcode, constants []object.BaseObject, messages []object.Message) *VM {
	return &VM{
		Instructions: bytecodes,
		Constants:    constants,
		Messages:     messages,
		Stack:        make([]object.BaseObject, StackSize),
		sp:           0,
		ip:           0,
		mp:           0,
		LastPopped:   nil,
		Env:          object.NewEnvironment(),
	}
}

func (vm *VM) pop() (object.BaseObject, object.BaseObject) {
	if vm.sp == 0 {
		return nil, object.NewError("stack underflow", 0)
	}
	o := vm.Stack[vm.sp-1]
	vm.sp--
	vm.LastPopped = o
	return o, nil
}

func (vm *VM) push(o object.BaseObject) object.BaseObject {
	if vm.sp >= StackSize {
		return object.NewError("stack overflow", o.Line())
	}

	vm.Stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) getMessage(p interface{}) interface{} {
	var val interface{}
	switch p.(type) {
	case int, int64:
		val = vm.Messages[vm.mp].(*object.IntMessage).Value
	case string:
		val = vm.Messages[vm.mp].(*object.StringMessage).Value
	default:
		panic(fmt.Sprintf("Cannot get message of type: %T", p))
	}

	vm.mp++

	return val
}

func (vm *VM) getIntMessage() int {
	val := vm.Messages[vm.mp].(*object.IntMessage).Value
	vm.mp++

	return val
}

func (vm *VM) getStringMessage() string {
	val := vm.Messages[vm.mp].(*object.StringMessage).Value
	vm.mp++

	return val
}

func (vm *VM) Run() object.BaseObject {
	for vm.ip < len(vm.Instructions) {
		bytecode := vm.Instructions[vm.ip]
		switch bytecode {
		case code.LoadConstant:
			index := vm.getIntMessage()
			err := vm.push(vm.Constants[index])
			if err != nil {
				return err
			}

		case code.Pop:
			_, e := vm.pop()
			if e != nil {
				return e
			}

		case code.Call:
			line := vm.getIntMessage()
			amountArgs := vm.getIntMessage()

			var arguments []object.BaseObject
			var v, e object.BaseObject
			for i := 0; i < amountArgs; i++ {
				v, e = vm.pop()
				if e != nil {
					return e
				}
				arguments = prependObj(arguments, v)
			}

			val, e := vm.pop()
			if e != nil {
				return e
			}

			var prt *object.Object
			if val.Parent() != nil {
				prt = val.Parent().(*object.Object)
			} else {
				prt = nil
			}
			ret := val.Call(vm.Env, prt, arguments, line)

			if isError(ret) {
				return ret
			}

			err := vm.push(ret)
			if err != nil {
				return err
			}

		case code.Method:
			name := vm.getStringMessage()

			val, e := vm.pop()
			if e != nil {
				return e
			}

			fn, er := val.FindMethod(name)
			if er != nil {
				return er
			}

			fn.SetParent(val)

			err := vm.push(fn)
			if err != nil {
				return err
			}

		case code.GetVar:
			name := vm.getStringMessage()
			line := vm.getIntMessage()

			v, ok := vm.Env.Get(name)
			if !ok {
				return object.NewError(fmt.Sprintf("name '%s' is not defined", name), line)
			}
			err := vm.push(v)
			if err != nil {
				return err
			}

		case code.SetVar:
			name := vm.getStringMessage()

			val, e := vm.pop()
			if e != nil {
				return e
			}

			vm.Env.Set(name, val)
			err := vm.push(val)
			if err != nil {
				return err
			}

		case code.Jump:
			index := vm.getIntMessage()
			mindex := vm.getIntMessage()
			vm.ip = index
			vm.mp = mindex

		case code.JumpFalse:
			index := vm.getIntMessage()
			mindex := vm.getIntMessage()
			val, e := vm.pop()
			if e != nil {
				return e
			}

			isTruthy, err := object.GetOneTruthy(val.(*object.Object), vm.Env, val.Line())
			if err != nil {
				return err
			}

			if !isTruthy {
				vm.ip = index
				vm.mp = mindex
			}

		default:
			return object.NewError(fmt.Sprintf("bytecode not supported: %d", bytecode), 0)
		}
		vm.ip++
	}
	return nil
}

func prependObj(x []object.BaseObject, y object.BaseObject) []object.BaseObject {
	x = append(x, nil)
	copy(x[1:], x)
	x[0] = y
	return x
}

func isError(obj object.BaseObject) bool {
	return obj.Type() == object.ERROR
}
