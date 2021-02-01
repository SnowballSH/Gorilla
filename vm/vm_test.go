package vm

import (
	"Gorilla/code"
	"Gorilla/object"
	"testing"
)

func assertStack(t *testing.T, vm *VM, stack []object.BaseObject) {
	err := vm.Run()
	if err != nil {
		t.Errorf("VM ERROR: %s", err)
		return
	}
	if vm.sp != len(stack) {
		t.Errorf("Stack length not same, expected length %d, got %d", len(stack), vm.sp)
		return
	}
	for i, v := range stack {
		if vm.Stack[i].Inspect() != v.Inspect() {
			t.Errorf("Stack Child %d not correct, expected %s, got %s", i, v.Inspect(), vm.Stack[i].Inspect())
		}
	}
}

func Test1(t *testing.T) {
	vm := New(
		[]code.Opcode{
			code.LoadConstant,
			code.LoadConstant,
			code.LoadConstant,
		},
		[]object.BaseObject{
			&object.Integer{
				Value: 123,
				SLine: 0,
			},
			&object.Integer{
				Value: 456,
				SLine: 0,
			},
		},
		[]object.Message{
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 1},
		},
	)
	stack := []object.BaseObject{
		&object.Integer{
			Value: 123,
			SLine: 0,
		},
		&object.Integer{
			Value: 123,
			SLine: 0,
		},
		&object.Integer{
			Value: 456,
			SLine: 0,
		},
	}
	assertStack(t, vm, stack)
}

func Test2(t *testing.T) {
	vm := New(
		[]code.Opcode{
			code.LoadConstant,
			code.LoadConstant,
			code.Pop,
			code.LoadConstant,
		},
		[]object.BaseObject{
			&object.Integer{
				Value: 123,
				SLine: 0,
			},
			&object.Integer{
				Value: 456,
				SLine: 0,
			},
		},
		[]object.Message{
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 1},
			&object.IntMessage{Value: 1},
		},
	)
	stack := []object.BaseObject{
		&object.Integer{
			Value: 123,
			SLine: 0,
		},
		&object.Integer{
			Value: 456,
			SLine: 0,
		},
	}
	assertStack(t, vm, stack)
}
