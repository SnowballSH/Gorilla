package token

type Token struct {
	Type    string
	Literal string
	Char    int // end char
	Line    int // end line
}

const (
	Integer = "Integer"
	Iden    = "Identifier"

	Semicolon = ";"
	Eq        = "="

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
