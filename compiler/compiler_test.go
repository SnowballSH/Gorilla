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
	assert.Equal(t, []byte{grammar.Magic, grammar.Integer, 2, 0xf8, 0x00, grammar.Pop}, comp)

	comp, ok = Compile(parser.NewParser(parser.NewLexer("69420\n\n624485")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 3, 0xAC, 0x9E, 0x04, grammar.Pop,
		grammar.Advance, grammar.Advance,
		grammar.Integer, 3, 0xE5, 0x8E, 0x26, grammar.Pop,
	}, comp)
}

func TestVars(t *testing.T) {
	comp, ok := Compile(parser.NewParser(parser.NewLexer("$ab12")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic, grammar.GetVar, 5, '$', 'a', 'b', '1', '2', grammar.Pop}, comp)

	comp, ok = Compile(parser.NewParser(parser.NewLexer("_X0 = \nY\n\n_X0")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.SetVar, 3, '_', 'X', '0', grammar.Advance, grammar.GetVar, 1, 'Y', grammar.Pop,
		grammar.Advance, grammar.Advance,
		grammar.GetVar, 3, '_', 'X', '0', grammar.Pop,
	}, comp)
}
