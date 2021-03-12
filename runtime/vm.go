package runtime

import (
	"ekyu.moe/leb128"
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/grammar"
)

type VM struct {
	source []byte
	ip     int

	line int

	stack []BaseObject

	error *errors.VMERROR

	lastPopped BaseObject
}

func NewVM(source []byte) *VM {
	return &VM{
		source: source,
		ip:     0,

		line: 0,

		stack: nil,

		error: nil,

		lastPopped: nil,
	}
}

func (vm *VM) Error(why string) {
	x := errors.MakeVMError(why, vm.line)
	vm.error = x
}

func (vm *VM) push(obj BaseObject) {
	vm.stack = append(vm.stack, obj)
}

func (vm *VM) pop() BaseObject {
	l := len(vm.stack) - 1
	k := vm.stack[l]
	vm.stack = vm.stack[:l]
	vm.lastPopped = k
	return k
}

func (vm *VM) read() byte {
	k := vm.source[vm.ip]
	vm.ip++
	return k
}

func (vm *VM) Run() {
	defer func() {
		if r := recover(); r != nil {
			errors.TestVMERR(r)
		}
	}()

	length := len(vm.source)

	if length == 0 || vm.read() != grammar.Magic {
		vm.Error("Not a valid Gorilla bytecode")
		panic(vm.error)
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
	case grammar.Pop:
		vm.pop()

	case grammar.Integer:
		length := int(vm.read())
		var number []byte
		for i := 0; i < length; i++ {
			number = append(number, vm.read())
		}
		val, _ := leb128.DecodeSleb128(number)
		vm.push(NewInteger(val))
	}
}
