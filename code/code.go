package code

type Opcode byte

// All Bytecodes
const (
	_ Opcode = iota

	// Loads a Constant
	LoadConstant

	// Pops the stack
	Pop

	// Call Function
	Call
	/*
		Syntax:
			Messages: [line int, amountArgs int]
			Stack:    [function, arg1, arg2, arg3, ...]
	*/

	// Call Method
	Method
	/*
		Syntax:
			Messages: [name string]
			Stack:    [callee]
	*/
)
