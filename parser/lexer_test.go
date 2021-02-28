package parser

import (
	"github.com/SnowballSH/Gorilla/parser/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLexer(t *testing.T) {
	lexer := NewLexer("你好, world!") // 10 chars
	assert.Equal(t, 10, lexer.inputLength)
}

func TestTerminator(t *testing.T) {
	lexer := NewLexer("\r\n \n")
	/*
		\r  \n  sp  \n
		0,0 1,0 0,1 1,1
	*/
	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "\r\n", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 1, n.Line)

	lexer = NewLexer("\n #a \r \n")
	/*
		\n  sp  #   a   sp  \n  \n
		0,0 0,1 1,1 2,1 3,1 4,1 0,2
	*/
	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 1, n.Char)
	assert.Equal(t, 0, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 5, n.Char)
	assert.Equal(t, 1, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 2, n.Line)

	assert.Equal(t, lexer.peek(66), rune(0))

	lexer = NewLexer("\r\n;")

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "\r\n", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Semicolon, n.Type)
	assert.Equal(t, 1, n.Char)
	assert.Equal(t, 1, n.Line)
}

func TestInteger(t *testing.T) {
	lexer := NewLexer("123 9 99")
	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Integer, n.Type)
	assert.Equal(t, 3, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "123", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Integer, n.Type)
	assert.Equal(t, 5, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "9", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Integer, n.Type)
	assert.Equal(t, 8, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "99", n.Literal)
}

func TestIden(t *testing.T) {
	lexer := NewLexer("var\n$global0 _hello123")
	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Iden, n.Type)
	assert.Equal(t, len("var"), n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "var", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, "\n", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Iden, n.Type)
	assert.Equal(t, len("$global0"), n.Char)
	assert.Equal(t, 1, n.Line)
	assert.Equal(t, "$global0", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Iden, n.Type)
	assert.Equal(t, 18, n.Char)
	assert.Equal(t, 1, n.Line)
	assert.Equal(t, "_hello123", n.Literal)
}

func TestMisc(t *testing.T) {
	lexer := NewLexer("a = 1")
	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Iden, n.Type)
	assert.Equal(t, len("a"), n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "a", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Eq, n.Type)
	assert.Equal(t, 3, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "=", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Integer, n.Type)
	assert.Equal(t, 5, n.Char)
	assert.Equal(t, 0, n.Line)
	assert.Equal(t, "1", n.Literal)
}

func TestBinOp(t *testing.T) {
	lexer := NewLexer("+-*/")
	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Plus, n.Type)
	assert.Equal(t, "+", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Minus, n.Type)
	assert.Equal(t, "-", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Star, n.Type)
	assert.Equal(t, "*", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.Slash, n.Type)
	assert.Equal(t, "/", n.Literal)

	n = lexer.next()
	assert.Equal(t, token.EOF, n.Type)
}

func TestIllegal(t *testing.T) {
	lexer := NewLexer("✔\rΣ")
	assert.Equal(t, 3, lexer.inputLength)

	var n token.Token

	n = lexer.next()
	assert.Equal(t, token.Illegal, n.Type)
	assert.Equal(t, 1, n.Char)
	assert.Equal(t, 0, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Newline, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 0, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Illegal, n.Type)
	assert.Equal(t, 1, n.Char)
	assert.Equal(t, 1, n.Line)
}
