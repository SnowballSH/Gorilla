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

	// Set Method
	SetMethod
	/*
		Syntax:
			Messages: [name string]
			Stack:    [receiver, value]
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

	// If Else
	Jump
	/*
		Syntax:
			Messages: [index int, message int]
			Stack:    []
	*/

	JumpFalse
	/*
		Syntax:
			Messages: [index int, message int]
			Stack:    [value]
	*/

	Return
	/*
		Syntax:
			Messages: []
			Stack:    [value]
	*/

	MakeArray
	/*
		Syntax:
			Messages: [amountValues int, line int]
			Stack:    [value1, value2, ...]
	*/

	MakeHash
	/*
		Syntax:
			Messages: [amountValues int, line int]
			Stack:    [key1, value1, key2, value2, ...]
	*/
)
