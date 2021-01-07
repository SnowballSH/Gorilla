package lexer

import (
	"testing"

	"../config"
	"../token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5
let ten = 10

let add = fn(x, y) {
  x + y
}

let result = add(five, ten)
!-/*5
5 < 10 > 5
<= >= <- ->

if (5 < 10) {
	return true
} else {
	return false
}

10 == 10
10 != 9
"hehehehe"
`

	nl := config.GetOSNewline("linux")

	tests := []struct {
		expectedType    token.TType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, nl},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, nl},
		{token.SEMICOLON, nl},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.SEMICOLON, nl},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, nl},
		{token.RBRACE, "}"},
		{token.SEMICOLON, nl},
		{token.SEMICOLON, nl},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, nl},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, nl},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, nl},
		{token.LTEQ, "<="},
		{token.GTEQ, ">="},
		{token.LARR, "<-"},
		{token.RARR, "->"},
		{token.SEMICOLON, nl},
		{token.SEMICOLON, nl},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.SEMICOLON, nl},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, nl},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.SEMICOLON, nl},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, nl},
		{token.RBRACE, "}"},
		{token.SEMICOLON, nl},
		{token.SEMICOLON, nl},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, nl},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, nl},
		{token.STRING, "hehehehe"},
		{token.SEMICOLON, nl},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
