package token

type TType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENTIFIER"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MOD      = "%"

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
	ELSIF        = "ELSIF"
	ELSE         = "ELSE"
	RETURN       = "RETURN"
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
	"elsif":  ELSIF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}