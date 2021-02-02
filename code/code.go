package code

type Opcode byte

// All Bytecodes
const (
	_ Opcode = iota

	// Loads a Constant
	LoadConstant

	// Pops the stack
	Pop

	// Call Method
	CallMethod
	/*
		Syntax:
			Messages: [line int, name string, amountArgs int]
			Stack:    [arg1, arg2, arg3, ..., callee]
	*/
)
