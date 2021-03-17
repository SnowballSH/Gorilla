package runtime

import (
	"ekyu.moe/leb128"
	"fmt"
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/grammar"
)

type VM struct {
	source []byte
	ip     int

	line int

	stack []BaseObject

	Error *errors.VMERROR

	LastPopped BaseObject
}

func NewVM(source []byte) *VM {
	return &VM{
		source: source,
		ip:     0,

		line: 0,

		stack: nil,

		Error: nil,

		LastPopped: nil,
	}
}

func (vm *VM) MakeError(why string) {
	x := errors.MakeVMError(why, vm.line)
	vm.Error = x
	panic(x)
}

func (vm *VM) push(obj BaseObject) {
	vm.stack = append(vm.stack, obj)
}

func (vm *VM) pop() BaseObject {
	l := len(vm.stack) - 1
	k := vm.stack[l]
	vm.stack = vm.stack[:l]
	vm.LastPopped = k
	return k
}

func (vm *VM) read() byte {
	k := vm.source[vm.ip]
	vm.ip++
	return k
}

func (vm *VM) readInt() int64 {
	length := int(vm.read())
	var number []byte
	for i := 0; i < length; i++ {
		number = append(number, vm.read())
	}
	val, _ := leb128.DecodeSleb128(number)
	return val
}

func (vm *VM) readString() string {
	length := int(vm.read())
	var bytes []byte
	for i := 0; i < length; i++ {
		bytes = append(bytes, vm.read())
	}
	return string(bytes)
}

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

func (vm *VM) RunStatement() {
	_type := vm.read()
	switch _type {
	case grammar.Advance:
		vm.line++
	case grammar.Back:
		vm.line--
	case grammar.Pop:
		vm.pop()

	case grammar.Integer:
		vm.push(NewInteger(vm.readInt()))

	case grammar.GetInstance:
		self := vm.pop()
		g := vm.readString()
		o, ok := self.InstanceVariableGet(g)
		if !ok {
			vm.MakeError(fmt.Sprintf("Attribute '%s' does not exist on '%s' (class '%s')", g, self.ToString(), self.Class().Name))
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
	}
}
