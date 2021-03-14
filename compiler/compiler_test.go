package compiler

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicNumbers(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("120")).Parse())
	assert.Equal(t, []byte{grammar.Magic, grammar.Integer, 2, 0xf8, 0x00, grammar.Pop}, compiler.result)

	compiler = NewCompiler()
	compiler.Compile(parser.NewParser(parser.NewLexer("69420\n\n624485")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 3, 0xAC, 0x9E, 0x04, grammar.Pop,
		grammar.Advance, grammar.Advance,
		grammar.Integer, 3, 0xE5, 0x8E, 0x26, grammar.Pop,
	}, compiler.result)
}

func TestInfix(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("1 * 3")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		1, 0x01,
		grammar.Integer, 1, 0x01,
		1, '*',
		grammar.GetInstance,
		grammar.Call,
		grammar.Pop,
	}, compiler.result)

	compiler = NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("2 + 1 * 3")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 0x03, // 3

		1, 0x01, // 1 arg

		grammar.Integer, 1, 0x01, // 1

		1, '*', // *

		grammar.GetInstance, // 1.*
		grammar.Call,        // 1.*(3)

		1, 0x01, // 1 arg

		grammar.Integer, 1, 0x02, // 2

		1, '+', // +

		grammar.GetInstance, // 2.+
		grammar.Call,        // 2.+(1.*(3))

		grammar.Pop,
	}, compiler.result)
}

func TestVars(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("$ab12")).Parse())
	assert.Equal(t, []byte{grammar.Magic, grammar.GetVar, 5, '$', 'a', 'b', '1', '2', grammar.Pop}, compiler.result)

	compiler = NewCompiler()
	compiler.Compile(parser.NewParser(parser.NewLexer("_X0 = \nY\n\n_X0")).Parse())
	assert.Equal(t, []byte{grammar.Magic, grammar.Advance,
		grammar.GetVar, 1, 'Y', grammar.Back, grammar.SetVar, 3, '_', 'X', '0', grammar.Pop,
		grammar.Advance, grammar.Advance, grammar.Advance,
		grammar.GetVar, 3, '_', 'X', '0', grammar.Pop,
	}, compiler.result)
}
