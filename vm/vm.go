package vm

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/code"
	"github.com/SnowballSH/Gorilla/object"
	"github.com/SnowballSH/Gorilla/stdlib"
)

type VM struct {
	Frame     *Frame
	LastFrame *Frame
}

func New(bytecodes []code.Opcode, constants []object.BaseObject, messages []object.Message) *VM {
	return &VM{
		Frame: NewFrame(bytecodes, constants, messages),
	}
}

func (vm *VM) pop() (object.BaseObject, object.BaseObject) {
	l := len(vm.Frame.Stack)
	if l == 0 {
		return nil, object.NewError("stack underflow", 0)
	}
	o := vm.Frame.Stack[l-1]
	vm.Frame.Stack = vm.Frame.Stack[:l-1]
	vm.Frame.LastPopped = o
	return o, nil
}

func (vm *VM) push(o object.BaseObject) {
	vm.Frame.Stack = append(vm.Frame.Stack, o)
}

func (vm *VM) getIntMessage() int {
	val := vm.Frame.Messages[vm.Frame.mp].(*object.IntMessage).Value
	vm.Frame.mp++

	return val
}

func (vm *VM) getStringMessage() string {
	val := vm.Frame.Messages[vm.Frame.mp].(*object.StringMessage).Value
	vm.Frame.mp++

	return val
}

func (vm *VM) Run() object.BaseObject {
	for vm.Frame.ip < len(vm.Frame.Instructions) {
		bytecode := vm.Frame.Instructions[vm.Frame.ip]
		switch bytecode {
		case code.LoadConstant:
			index := vm.getIntMessage()
			obj := vm.Frame.Constants[index]

			if obj.Type() == object.FUNCTION {
				obj.(*object.Object).InternalValue.(*object.FunctionValue).FreeEnv = vm.Frame.Env
				obj.(*object.Object).CallFunc = RunFunc(obj.(*object.Object), vm)
			} else if obj.Type() == object.MACRO {
				obj.(*object.Object).InternalValue.(*object.FunctionValue).FreeEnv = vm.Frame.Env
				obj.(*object.Object).CallFunc = RunDo(obj.(*object.Object), vm)
			}

			vm.push(obj)

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

			// Call function
			vv := val.(*object.Object)

			ret := vv.Call(vm.Frame.Env, prt, arguments, line)

			if IsError(ret) {
				return ret
			}

			vm.push(ret)

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

			fn.SetParent(val.(*object.Object))

			vm.push(fn)

		case code.SetMethod:
			name := vm.getStringMessage()

			val, e := vm.pop()
			if e != nil {
				return e
			}

			rec, e := vm.pop()
			if e != nil {
				return e
			}

			rec.SetMethod(name, val)
			vm.push(rec)

		case code.GetVar:
			name := vm.getStringMessage()
			line := vm.getIntMessage()

			v, ok := vm.Frame.Env.Get(name)
			if !ok {
				return object.NewError(fmt.Sprintf("name '%s' is not defined", name), line)
			}
			vm.push(v)

		case code.SetVar:
			name := vm.getStringMessage()

			val, e := vm.pop()
			if e != nil {
				return e
			}

			vm.Frame.Env.Set(name, val, vm.Frame.Macro)
			vm.push(val)

		case code.Jump:
			index := vm.getIntMessage()
			mindex := vm.getIntMessage()
			vm.Frame.ip = index
			vm.Frame.mp = mindex

		case code.JumpFalse:
			index := vm.getIntMessage()
			mindex := vm.getIntMessage()
			val, e := vm.pop()
			if e != nil {
				return e
			}

			isTruthy, err := object.GetOneTruthy(val.(*object.Object), vm.Frame.Env, val.Line())
			if err != nil {
				return err
			}

			if !isTruthy {
				vm.Frame.ip = index
				vm.Frame.mp = mindex
			}

		case code.Return:
			_, e := vm.pop()
			if e != nil {
				return e
			}

			vm.Frame.ip = len(vm.Frame.Instructions)

		case code.MakeArray:
			amountVals := vm.getIntMessage()
			line := vm.getIntMessage()
			var values []object.BaseObject
			for i := 0; i < amountVals; i++ {
				val, e := vm.pop()
				if e != nil {
					return e
				}
				values = prependObj(values, val)
			}
			vm.push(object.NewArray(values, line))

		case code.MakeHash:
			amountVals := vm.getIntMessage()
			line := vm.getIntMessage()
			pairs := map[object.HashKey]*object.HashValue{}
			for i := 0; i < amountVals; i++ {
				value, e := vm.pop()
				if e != nil {
					return e
				}
				key, e := vm.pop()
				if e != nil {
					return e
				}

				hashedKey, ok := object.HashObject(key)
				if !ok {
					return object.NewError(fmt.Sprintf("Type '%s' is not hashable", key.Type()), line)
				}

				pairs[hashedKey] = &object.HashValue{
					Key:   key,
					Value: value,
				}
			}
			vm.push(object.NewHash(pairs, line))

		case code.Import:
			name := vm.getStringMessage()
			line := vm.getIntMessage()
			lib, ok := stdlib.StandardLibrary[name]
			if !ok {
				return object.NewError(fmt.Sprintf("Module not found: %s", name), line)
			}
			object.GlobalBuiltins[name] = lib

		default:
			return object.NewError(fmt.Sprintf("bytecode not supported: %d", bytecode), 0)
		}
		vm.Frame.ip++

		/*
			println("[")
			for _, o := range vm.Stack[:vm.sp] {
				if o != nil {
					if o.Parent() != nil {
						println("PARENT: " + o.Parent().Debug(), o.Parent().(*object.Object))
					}
					println(o.Debug(), o.(*object.Object))
				}
			}
			println("]")
		*/
	}
	return nil
}

func prependObj(x []object.BaseObject, y object.BaseObject) []object.BaseObject {
	x = append(x, nil)
	copy(x[1:], x)
	x[0] = y
	return x
}

func IsError(obj object.BaseObject) bool {
	return obj.Type() == object.ERROR
}

func RunFunc(vv *object.Object, vmm *VM) func(env *object.Environment, self *object.Object, args []object.BaseObject, line int) object.BaseObject {
	if vv.Type() != object.FUNCTION {
		return vv.Call
	}

	return func(env *object.Environment, self *object.Object, args []object.BaseObject, line int) object.BaseObject {
		fstr := vv.Value().(*object.FunctionValue)

		if len(args) != len(fstr.Params) {
			return object.NewError(
				fmt.Sprintf("Argument amount mismatch: Expected %d, got %d", len(fstr.Params), len(args)),
				line,
			)
		}

		newframe := NewFrame(fstr.Bytecodes, fstr.Constants, fstr.Messages)
		newframe.Env = object.NewEnclosedEnvironment(env)

		for name, free := range fstr.FreeEnv.Store {
			newframe.Env.Set(name, free, false)
		}

		for i, vvv := range fstr.Params {
			newframe.Env.Set(vvv, args[i], false)
		}

		newframe.LastFrame = vmm.Frame
		vmm.Frame = newframe

		e := vmm.Run()
		if e != nil {
			return e
		}

		last := vmm.Frame.LastPopped

		vmm.Frame = vmm.Frame.LastFrame

		if last == nil {
			return object.NULLOBJ
		}

		if IsError(last) {
			return last
		}
		return last
	}
}

func RunDo(vv *object.Object, vmm *VM) func(env *object.Environment, self *object.Object, args []object.BaseObject, line int) object.BaseObject {
	if vv.Type() != object.MACRO {
		return vv.Call
	}

	return func(env *object.Environment, self *object.Object, args []object.BaseObject, line int) object.BaseObject {
		fstr := vv.Value().(*object.FunctionValue)

		if len(args) != len(fstr.Params) {
			return object.NewError(
				fmt.Sprintf("Argument amount mismatch: Expected %d, got %d", len(fstr.Params), len(args)),
				line,
			)
		}

		newframe := NewFrame(fstr.Bytecodes, fstr.Constants, fstr.Messages)
		newframe.Macro = true
		newframe.Env = env

		for i, vvv := range fstr.Params {
			newframe.Env.Free[vvv] = args[i]
		}

		newframe.LastFrame = vmm.Frame
		vmm.Frame = newframe

		e := vmm.Run()
		if e != nil {
			return e
		}

		last := vmm.Frame.LastPopped

		vmm.Frame = vmm.Frame.LastFrame

		vmm.Frame.Env.Free = map[string]object.BaseObject{}

		if last == nil {
			return object.NULLOBJ
		}

		if IsError(last) {
			return last
		}
		return last
	}
}
