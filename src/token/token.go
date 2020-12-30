package token

type TType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENTIFIER" // add, foobar, x, y, ...
	INT   = "INT"        // 1343456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT   = "<"
	GT   = ">"
	LTEQ = "<="
	GTEQ = ">="

	EQ  = "=="
	NEQ = "!="

	LARR = "<-"
	RARR = "->"

	// Delimiters
	COMMA     = ","
	SEMICOLON = "; or newline"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION     = "FUNCTION"
	STMTFUNCTION = "STMTFUNCTION"
	LET          = "LET"
	TRUE         = "TRUE"
	FALSE        = "FALSE"
	IF           = "IF"
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
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
