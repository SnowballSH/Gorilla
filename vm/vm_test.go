package vm

import (
	"Gorilla/code"
	"Gorilla/object"
	"testing"
)

func assertStack(t *testing.T, vm *VM, stack []object.BaseObject, result object.BaseObject) {
	err := vm.Run()
	if err != nil && isError(err) {
		t.Errorf("VM ERROR: %s", err)
		return
	}
	if result != vm.lastPopped {
		var val1, val2 string
		if result != nil {
			val1 = result.Inspect()
		} else {
			val1 = "nil"
		}
		if vm.lastPopped != nil {
			val2 = vm.lastPopped.Inspect()
		} else {
			val2 = "nil"
		}

		if val1 == val2 {
		} else {
			t.Errorf("Result not match: expected %s, got %s", val1, val2)
			return
		}
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

var testInt1 = object.NewInteger(
	123,
	0,
)

var testInt2 = object.NewInteger(
	456,
	0,
)

func Test1(t *testing.T) {
	vm := New(
		[]code.Opcode{
			code.LoadConstant,
			code.LoadConstant,
			code.LoadConstant,
		},
		[]object.BaseObject{
			testInt1,
			testInt2,
		},
		[]object.Message{
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 1},
		},
	)
	stack := []object.BaseObject{
		testInt1,
		testInt1,
		testInt2,
	}
	assertStack(t, vm, stack, nil)
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
			testInt1,
			testInt2,
		},
		[]object.Message{
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 1},
			&object.IntMessage{Value: 1},
		},
	)
	stack := []object.BaseObject{
		testInt1,
		testInt2,
	}
	assertStack(t, vm, stack, testInt2)
}

func Test3(t *testing.T) {
	vm := New(
		[]code.Opcode{
			code.LoadConstant,
			code.LoadConstant,
			code.CallMethod,
			code.Pop,
		},
		[]object.BaseObject{
			object.NewInteger(2, 0),
			object.NewInteger(1, 0),
		},
		[]object.Message{
			&object.IntMessage{Value: 0},
			&object.IntMessage{Value: 1},
			&object.IntMessage{Value: 0},
			&object.StringMessage{Value: "add"},
			&object.IntMessage{Value: 1},
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NewInteger(3, 0))
}
