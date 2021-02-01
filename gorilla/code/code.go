package code

type opcode byte

// All Bytecodes
const (
	_ opcode = iota

	// Loads a Constant
	LoadConstant

	// Pops the stack
	Pop

	// Binary Operations
	Addition
	Subtract
	Multiply
	Division
)
