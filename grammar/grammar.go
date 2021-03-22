package grammar

const (
	Magic byte = 0x69

	Pop byte = iota
	Advance
	Back

	Integer
	// Length of integer in unsigned leb128
	// Followed by leb128

	String
	// String in bytes

	GetVar // name string

	SetVar // name string
	// Expression on stack

	GetInstance // name string
	// Expression on stack

	Call // numberOfArgs int
	// [arg3, arg2, arg1, callee]
	// Will perform callee(arg1, arg2, arg3)
)
