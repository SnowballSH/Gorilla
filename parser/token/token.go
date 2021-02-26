package token

type Token struct {
	Type    string
	Literal string
	Char    int // end char
	Line    int // end line
}

const (
	Integer    = "Integer"
	Terminator = "Terminator"
)
