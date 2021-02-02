package vm

import (
	"Gorilla/code"
	"Gorilla/object"
	"testing"
)

func assertStack(t *testing.T, vm *VM, stack []object.BaseObject, result object.BaseObject) {
	err := vm.Run()
	if err != nil && isError(err) {
		t.Errorf("VM ERROR: %s", err.Inspect())
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
			object.NewMessage(0),
			object.NewMessage(0),
			object.NewMessage(1),
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
			object.NewMessage(0),
			object.NewMessage(1),
			object.NewMessage(1),
		},
	)
	stack := []object.BaseObject{
		testInt1,
		testInt2,
	}
	assertStack(t, vm, stack, testInt2)
}

func Test3(t *testing.T) {
	// Pseudo Code: `2.add(1)` or `2 + 1`
	vm := New(
		[]code.Opcode{
			code.LoadConstant, // Load 2
			code.Method,       // Find 2's "add"
			code.LoadConstant, // Load 1
			code.Call,         // Call 2's"add"
			code.Pop,
		},
		[]object.BaseObject{
			object.NewInteger(2, 0), // Constant: 2
			object.NewInteger(1, 0), // Constant: 1
		},
		[]object.Message{
			object.NewMessage(0),     // Load 2
			object.NewMessage("add"), // Add
			object.NewMessage(1),     // Load 1
			object.NewMessage(66),    // Line 66
			object.NewMessage(1),     // 1 argument
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NewInteger(3, 66))
}

func Test4(t *testing.T) {
	// Pseudo Code: `a = 4; a`
	vm := New(
		[]code.Opcode{
			code.LoadConstant, // Load 4
			code.SetVar,       // Set a to 4
			code.Pop,
			code.GetVar, // Load a
			code.Pop,
		},
		[]object.BaseObject{
			object.NewInteger(4, 0), // Constant: 4
		},
		[]object.Message{
			object.NewMessage(0),
			object.NewMessage("a"),
			object.NewMessage("a"),
			object.NewMessage(66),
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NewInteger(4, 66))
}
