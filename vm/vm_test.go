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
	if result != vm.LastPopped {
		var val1, val2 string
		if result != nil {
			val1 = result.Inspect()
		} else {
			val1 = "nil"
		}
		if vm.LastPopped != nil {
			val2 = vm.LastPopped.Inspect()
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

func Test5(t *testing.T) {
	// Pseudo Code: `if true {1} else {2}`
	vm := New(
		[]code.Opcode{
			code.LoadConstant, // Load true                     0
			code.JumpFalse,    // Jump If true is false         1
			code.LoadConstant, // Load 1                        2
			code.Jump,         // Jump Out                      3
			code.LoadConstant, // Else, Load 2                  4
			code.Pop,          // Pop                           5
		},
		[]object.BaseObject{
			object.NewBool(true, 0),
			object.NewInteger(1, 0),
			object.NewInteger(2, 0),
		},
		[]object.Message{
			object.NewMessage(0), // Load true             0
			object.NewMessage(3), // Jump to 3 -> 4        1
			object.NewMessage(6), // Jump Message to 6     2
			object.NewMessage(1), // Load 1                3
			object.NewMessage(4), // Jump to 4 -> 5        4
			object.NewMessage(7), // Jump Message to 7     5
			object.NewMessage(2), // Load 2                6
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NewInteger(1, 0))
}

func Test6(t *testing.T) {
	// Pseudo Code: `if false {1} else {2}`
	vm := New(
		[]code.Opcode{
			code.LoadConstant, // Load false                    0
			code.JumpFalse,    // Jump If false is false        1
			code.LoadConstant, // Load 1                        2
			code.Jump,         // Jump Out                      3
			code.LoadConstant, // Else, Load 2                  4
			code.Pop,          // Pop                           5
		},
		[]object.BaseObject{
			object.NewBool(false, 0),
			object.NewInteger(1, 0),
			object.NewInteger(2, 0),
		},
		[]object.Message{
			object.NewMessage(0), // Load false             0
			object.NewMessage(3), // Jump to 3 -> 4        1
			object.NewMessage(6), // Jump Message to 6     2
			object.NewMessage(1), // Load 1                3
			object.NewMessage(4), // Jump to 4 -> 5        4
			object.NewMessage(7), // Jump Message to 7     5
			object.NewMessage(2), // Load 2                6
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NewInteger(2, 0))
}

func Test7(t *testing.T) {
	// Pseudo Code: `i = 5; while i { i = i - 1 }`
	vm := New(
		[]code.Opcode{
			code.LoadConstant, // Load 5                        0
			code.SetVar,       // Set i to 5                    1
			code.Pop,          // Pop                           2
			code.GetVar,       // Get Variable 'i'              3
			code.JumpFalse,    // Jump                          4
			code.GetVar,       // Get Variable 'i'              5
			code.Method,       // Find i's "sub"                6
			code.LoadConstant, // Load 1                        7
			code.Call,         // Call i's "sub"                8
			code.SetVar,       // Set i to i - 1                9
			code.Pop,          // Pop                           10
			code.Jump,         // Jump Back                     11
			code.LoadConstant, // Load NULL                     12
			code.Pop,          // Pop                           13
		},
		[]object.BaseObject{
			object.NewInteger(5, 0),
			object.NewInteger(1, 0),
			object.NULLOBJ,
		},
		[]object.Message{
			object.NewMessage(0),     // Load 5            0
			object.NewMessage("i"),   //                   1
			object.NewMessage("i"),   //                   2
			object.NewMessage(0),     //                   3
			object.NewMessage(11),    // 11 -> 12          4
			object.NewMessage(15),    // 15                5
			object.NewMessage("i"),   //                   6
			object.NewMessage(0),     //                   7
			object.NewMessage("sub"), //                   8
			object.NewMessage(1),     //                   9
			object.NewMessage(0),     //                   10
			object.NewMessage(1),     //                   11
			object.NewMessage("i"),   //                   12
			object.NewMessage(2),     // 2 -> 3            13
			object.NewMessage(2),     // 2                 14
			object.NewMessage(2),     // NULL              15
		},
	)
	var stack []object.BaseObject
	assertStack(t, vm, stack, object.NULLOBJ)
}
