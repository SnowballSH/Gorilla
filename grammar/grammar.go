package grammar

const (
	Magic byte = 0x69

	Pop byte = iota
	Advance
	Back

	Integer
	// Length of integer in unsigned leb128
	// Followed by leb128

	GetVar
	// Length of name (string) in bytes
	// String encoded in bytes

	SetVar
	// Expression on stack
	// Length of name (string) in bytes
	// String encoded in bytes

	GetInstance
	// Expression on stack
	// Length of name (string) in bytes
	// String encoded in bytes

	Call
	// stack should be [arg3, arg2, arg1, numberOfArgs, callee]
	// Will perform callee(arg1, arg2, arg3)
)
