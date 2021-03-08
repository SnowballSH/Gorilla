package grammar

const (
	Magic byte = 0x69

	Integer byte = iota
	// Length of integer in unsigned leb128
	// Followed by leb128

	GetVar
	// Length of name (string) in bytes
	// String encoded in bytes

	SetVar
	// Length of name (string) in bytes
	// String encoded in bytes
	// Expression
)
