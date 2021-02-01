package vm

import (
	"Gorilla/code"
	"Gorilla/object"
	"fmt"

	_ "github.com/alecthomas/participle"
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
	}
}

func (vm *VM) pop() (object.BaseObject, error) {
	if vm.sp == 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	o := vm.Stack[vm.sp-1]
	vm.sp--
	return o, nil
}

func (vm *VM) push(o object.BaseObject) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.Stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) Run() error {
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

		case code.Addition:

		default:
			return fmt.Errorf("bytecode not supported: %d", bytecode)
		}
		vm.ip++
	}
	return nil
}
