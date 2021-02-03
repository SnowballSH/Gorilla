package token

type TType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENTIFIER"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

	POW = "**"

	OR  = "||"
	AND = "&&"

	BOR  = "|"
	XOR  = "^"
	BAND = "&"

	LT   = "<"
	GT   = ">"
	LTEQ = "<="
	GTEQ = ">="

	DOT = "."

	EQ  = "=="
	NEQ = "!="

	LARR = "<-"
	RARR = "->"

	// Delimiters
	COMMA     = ","
	SEMICOLON = "; or newline"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION     = "FUNCTION"
	STMTFUNCTION = "STMTFUNCTION"
	LET          = "LET"
	TRUE         = "TRUE"
	FALSE        = "FALSE"
	IF           = "IF"
	ELSE         = "ELSE"
	RETURN       = "RETURN"

	WHILE = "WHILE"
)

type Token struct {
	Type    TType
	Literal string
	Line    int
}

var keywords = map[string]TType{
	"fn":     FUNCTION,
	"func":   STMTFUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
}

func LookupIdent(ident string) TType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}