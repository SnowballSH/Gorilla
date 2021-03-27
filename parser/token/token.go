package token

type Token struct {
	Type    string
	Literal string
	Char    int // end char
	Line    int // end line
}

const (
	Integer = "Integer"
	String  = "String"

	Iden = "Identifier"

	Semicolon = ";"
	Eq        = "="
	Comma     = ","

	Plus  = "+"
	Minus = "-"
	Star  = "*"
	Slash = "/"

	LParen = "("
	RParen = ")"

	Newline = "Newline"
	Illegal = "Illegal"
	EOF     = "End of file"
)
