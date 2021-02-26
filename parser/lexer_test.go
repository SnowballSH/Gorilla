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
	assert.Equal(t, token.Terminator, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 0, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Terminator, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 1, n.Line)

	lexer = NewLexer("\n #a \r \n")
	/*
		\n  sp  #   a   sp  \n  \n
		0,0 0,1 1,1 2,1 3,1 4,1 0,2
	*/
	n = lexer.next()
	assert.Equal(t, token.Terminator, n.Type)
	assert.Equal(t, 1, n.Char)
	assert.Equal(t, 0, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Terminator, n.Type)
	assert.Equal(t, 5, n.Char)
	assert.Equal(t, 1, n.Line)

	n = lexer.next()
	assert.Equal(t, token.Terminator, n.Type)
	assert.Equal(t, 2, n.Char)
	assert.Equal(t, 2, n.Line)

	assert.Equal(t, lexer.peek(66), rune(0))
}
