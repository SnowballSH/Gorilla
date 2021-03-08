package compiler

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicNumbers(t *testing.T) {
	comp, ok := Compile(parser.NewParser(parser.NewLexer("120")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic, grammar.Integer, 2, 0xf8, 0x00}, comp)

	comp, ok = Compile(parser.NewParser(parser.NewLexer("69420")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic, grammar.Integer, 3, 0xac, 0x9e, 0x04}, comp)
}

func TestVars(t *testing.T) {
	comp, ok := Compile(parser.NewParser(parser.NewLexer("$ab12")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic, grammar.GetVar, 5, '$', 'a', 'b', '1', '2'}, comp)

	comp, ok = Compile(parser.NewParser(parser.NewLexer("_X0 = Y")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic, grammar.SetVar, 3, '_', 'X', '0', grammar.GetVar, 1, 'Y'}, comp)
}
