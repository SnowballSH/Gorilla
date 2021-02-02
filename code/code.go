package code

type Opcode byte

// All Bytecodes
const (
	_ Opcode = iota

	// Loads a Constant
	LoadConstant
	/*
		Syntax:
			Messages: [index int]
			Stack:    []
	*/

	// Pops the stack
	Pop
	/*
		Syntax:
			Messages: []
			Stack:    [(at least 1 object)]
	*/

	// Call Function
	Call
	/*
		Syntax:
			Messages: [line int, amountArgs int]
			Stack:    [function, arg1, arg2, arg3, ...]
	*/

	// Get Method
	Method
	/*
		Syntax:
			Messages: [name string]
			Stack:    [callee]
	*/

	// Get Var
	GetVar
	/*
		Syntax:
			Messages: [name string line int]
			Stack:    []
	*/

	// Set Var
	SetVar
	/*
		Syntax:
			Messages: [name string]
			Stack:    [value]
	*/
)
