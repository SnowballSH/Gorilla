package grammar

const (
	// The Magic Number
	Magic byte = 0x69

	// Pops the stack
	Pop byte = iota

	// Do nothing
	Noop

	// line++
	Advance

	// line--
	Back

	// Length of integer in unsigned leb128;
	// Followed by leb128
	Integer

	// String in bytes
	String

	// Null
	Null

	// name string
	GetVar

	// name string;
	// Expression on stack
	SetVar

	// name string;
	// Expression on stack
	GetInstance

	// numberOfArgs int;
	// [arg3, arg2, arg1, callee];
	// Will perform callee(arg1, arg2, arg3)
	Call

	// where int;
	// Expression on stack
	JumpIfFalse

	// where int
	Jump

	// numberOfParams int;
	// [param3, param2, param1, ...]
	// lengthOfCode int;
	// code
	Lambda

	// lengthOfCode int;
	// code
	Closure
)
