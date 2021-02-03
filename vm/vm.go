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

	env *object.Environment
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
		env:          object.NewEnvironment(),
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

func (vm *VM) Run() object.BaseObject {
	for vm.ip < len(vm.Instructions) {
		bytecode := vm.Instructions[vm.ip]
		switch bytecode {
		case code.LoadConstant:
			index := vm.Messages[vm.mp].(*object.IntMessage).Value
			vm.mp++
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
			line := vm.Messages[vm.mp].(*object.IntMessage).Value
			vm.mp++
			amountArgs := vm.Messages[vm.mp].(*object.IntMessage).Value
			vm.mp++

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

			ret := val.Call(vm.env, val.Parent().(*object.Object), arguments, line)

			if isError(ret) {
				return ret
			}

			err := vm.push(ret)
			if err != nil {
				return err
			}

		case code.Method:
			name := vm.Messages[vm.mp].(*object.StringMessage).Value
			vm.mp++

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
			name := vm.Messages[vm.mp].(*object.StringMessage).Value
			vm.mp++
			line := vm.Messages[vm.mp].(*object.IntMessage).Value
			vm.mp++

			v, ok := vm.env.Get(name)
			if !ok {
				return object.NewError(fmt.Sprintf("name '%s' is not defined", name), line)
			}
			err := vm.push(v)
			if err != nil {
				return err
			}

		case code.SetVar:
			name := vm.Messages[vm.mp].(*object.StringMessage).Value
			vm.mp++

			val, e := vm.pop()
			if e != nil {
				return e
			}

			vm.env.Set(name, val)
			err := vm.push(val)
			if err != nil {
				return err
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
