package code

type Opcode byte

// All Bytecodes
const (
	_ Opcode = iota

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
