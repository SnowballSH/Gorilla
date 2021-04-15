package runtime

import (
	"ekyu.moe/leb128"
	"fmt"
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/grammar"
)

// The Base VM struct
type VM struct {
	source []byte
	ip     int

	line int

	stack []BaseObject

	Error *errors.VMERROR

	LastPopped BaseObject

	Environment *Environment
}

// NewVM creates a new vm from array of bytes
func NewVM(source []byte) *VM {
	return &VM{
		source: source,
		ip:     0,

		line: 0,

		stack: nil,

		Error: nil,

		LastPopped: nil,

		Environment: NewEnvironment(),
	}
}

// NewVMWithStore is NewVM but with set environment
func NewVMWithStore(source []byte, env *Environment) *VM {
	return &VM{
		source: source,
		ip:     0,

		line: 0,

		stack: nil,

		Error: nil,

		LastPopped: nil,

		Environment: NewEnvironmentWithStore(env.Store),
	}
}

// MakeError panics a VM error
func (vm *VM) MakeError(why string) {
	x := errors.MakeVMError(why, vm.line)
	vm.Error = x
	panic(x)
}

// push pushes object to the end of stack
func (vm *VM) push(obj BaseObject) {
	vm.stack = append(vm.stack, obj)
}

// pop pops an object from the end of stack
func (vm *VM) pop() BaseObject {
	l := len(vm.stack) - 1
	k := vm.stack[l]
	vm.stack = vm.stack[:l]
	vm.LastPopped = k
	return k
}

// read reads an instruction
func (vm *VM) read() byte {
	k := vm.source[vm.ip]
	vm.ip++
	return k
}

// readInt reads an integer
func (vm *VM) readInt() int64 {
	length := int(vm.read())
	var number []byte
	for i := 0; i < length; i++ {
		number = append(number, vm.read())
	}
	val, _ := leb128.DecodeSleb128(number)
	return val
}

// readString reads a string
func (vm *VM) readString() string {
	length := int(vm.read())
	var bytes []byte
	for i := 0; i < length; i++ {
		bytes = append(bytes, vm.read())
	}
	return string(bytes)
}

// Run runs the bytecode
func (vm *VM) Run() {
	defer func() {
		if r := recover(); r != nil {
			errors.TestVMERR(r)
		}
	}()

	length := len(vm.source)

	if length == 0 || vm.read() != grammar.Magic {
		vm.MakeError("Not a valid Gorilla bytecode")
	}

	for vm.ip < length {
		vm.RunStatement()
	}
}

// RunStatement runs a single statement/opcode
func (vm *VM) RunStatement() {
	_type := vm.read()
	switch _type {
	case grammar.Advance:
		vm.line++
	case grammar.Back:
		vm.line--
	case grammar.Pop:
		vm.pop()
	case grammar.Noop:
		// Do nothing

	case grammar.Null:
		vm.push(Null)

	case grammar.Integer:
		vm.push(NewInteger(vm.readInt()))

	case grammar.String:
		vm.push(NewString(vm.readString()))

	case grammar.GetVar:
		name := vm.readString()
		o, ok := vm.Environment.Get(name)
		if !ok {
			o, ok = Global.Get(name)
		}
		if !ok {
			vm.MakeError(fmt.Sprintf("Variable '%s' is not defined", name))
		}
		vm.push(o)

	case grammar.SetVar:
		name := vm.readString()
		value := vm.pop()

		for _, w := range ReservedKW {
			if w == name {
				vm.MakeError(fmt.Sprintf("Variable name '%s' is reserved", name))
			}
		}

		if len(name) >= 1 && name[0] == '$' {
			Global.Set(name, value)
		} else {
			vm.Environment.Set(name, value)
		}
		vm.push(value)

	case grammar.GetInstance:
		self := vm.pop()
		g := vm.readString()
		o, ok := self.InstanceVariableGet(g)
		if !ok {
			vm.MakeError(fmt.Sprintf("Attribute '%s' does not exist on '%s' (class '%s')", g, self.Inspect(), self.Class().Name))
		}
		o.SetParent(self)
		vm.push(o)

	case grammar.Call:
		amount := vm.readInt()

		o := vm.pop()

		var args []BaseObject
		for i := int64(0); i < amount; i++ {
			args = append(args, vm.pop())
		}

		val, err := o.Call(o, args...)

		if err != nil {
			vm.MakeError(err.Error())
		}

		vm.push(val)

	case grammar.Lambda:
		amount := vm.readInt()
		var args []string
		for i := int64(0); i < amount; i++ {
			args = append(args, vm.readString())
		}

		var bc []byte
		length := vm.readInt()
		for i := int64(0); i < length; i++ {
			bc = append(bc, vm.read())
		}

		vm.push(NewLambda(args, bc, vm))

	case grammar.Jump:
		where := vm.readInt()
		vm.ip = int(where) + 1

	case grammar.JumpIfFalse:
		where := vm.readInt()
		if !vm.pop().IsTruthy() {
			vm.ip = int(where) + 1
		}
	}
}
