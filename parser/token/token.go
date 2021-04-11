package token

// The base Token struct
type Token struct {
	Type    string
	Literal string
	Char    int // end char
	Line    int // end line
}

const (
	Integer = "Integer"
	String  = "String"

	If   = "if"
	Else = "else"

	Iden = "Identifier"

	Semicolon = ";"
	Eq        = "="
	Comma     = ","

	DbEq = "=="
	Neq  = "!="
	Not  = "!"

	Plus    = "+"
	Minus   = "-"
	Star    = "*"
	Slash   = "/"
	Percent = "%"

	LParen = "("
	RParen = ")"

	LCurly = "{"
	RCurly = "}"

	VBar = "|"

	Newline = "Newline"
	Illegal = "Illegal"
	EOF     = "End of file"
)

// All Keywords
var Keywords = map[string]string{
	"if":   If,
	"else": Else,
}
