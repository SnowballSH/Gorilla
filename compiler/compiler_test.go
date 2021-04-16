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
	assert.Equal(t, []byte{grammar.Magic, grammar.Integer, 2, 0xf8, 0x00, grammar.Pop}, compiler.Result)

	compiler = NewCompiler()
	compiler.Compile(parser.NewParser(parser.NewLexer("69420\n\n624485")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 3, 0xAC, 0x9E, 0x04, grammar.Pop,
		grammar.Advance, grammar.Advance,
		grammar.Integer, 3, 0xE5, 0x8E, 0x26, grammar.Pop,
	}, compiler.Result)
}

func TestBasicString(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("'AB\nCD'")).Parse())
	assert.Equal(t, []byte{
		grammar.Magic,
		grammar.Advance,
		grammar.Back,
		grammar.String, 5, 'A', 'B', '\n', 'C', 'D', grammar.Pop,
	}, compiler.Result)
}

func TestInfix(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("1 * 3")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '*',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	}, compiler.Result)

	compiler = NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("2 + 1 * 3")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 0x03, // 3

		grammar.Integer, 1, 0x01, // 1

		grammar.GetInstance, // 1.*
		1, '*',              // *

		grammar.Call, // 1.*(3)
		1, 0x01,      // 1 arg

		grammar.Integer, 1, 0x02, // 2

		grammar.GetInstance, // 2.+
		1, '+',              // +

		grammar.Call, // 2.+(1.*(3))
		1, 0x01,      // 1 arg

		grammar.Pop,
	}, compiler.Result)

	compiler = NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("-1")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		2, '-', '@',
		grammar.Call,
		1, 0x0,
		grammar.Pop,
	}, compiler.Result)
}

func TestVars(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("$ab12")).Parse())
	assert.Equal(t, []byte{grammar.Magic, grammar.GetVar, 5, '$', 'a', 'b', '1', '2', grammar.Pop}, compiler.Result)

	compiler = NewCompiler()
	compiler.Compile(parser.NewParser(parser.NewLexer("_X0 = \nY\n\n_X0")).Parse())
	assert.Equal(t, []byte{grammar.Magic, grammar.Advance,
		grammar.GetVar, 1, 'Y', grammar.Back, grammar.SetVar, 3, '_', 'X', '0', grammar.Pop,
		grammar.Advance, grammar.Advance, grammar.Advance,
		grammar.GetVar, 3, '_', 'X', '0', grammar.Pop,
	}, compiler.Result)
}

func TestCall(t *testing.T) {
	compiler := NewCompiler()

	compiler.Compile(parser.NewParser(parser.NewLexer("$ab12(1, 2)")).Parse())
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 2,
		grammar.Integer, 1, 1,
		grammar.GetVar, 5, '$', 'a', 'b', '1', '2',
		grammar.Call, 1, 2,
		grammar.Pop}, compiler.Result)
}

func TestIfElse(t *testing.T) {
	compiler := NewCompiler()

	p := parser.NewParser(parser.NewLexer(`if 1 {

} else {5
}
1`))
	res := p.Parse()

	compiler.Compile(res)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 1,
		grammar.JumpIfFalse, 1, 10,
		grammar.Null,
		grammar.Jump, 1, 18,
		grammar.Advance,
		grammar.Advance,
		grammar.Integer, 1, 5, grammar.Noop,
		grammar.Back,
		grammar.Back,
		grammar.Pop,
		grammar.Advance, grammar.Advance, grammar.Advance, grammar.Advance,
		grammar.Integer, 1, 1, grammar.Pop,
	}, compiler.Result)

	compiler = NewCompiler()

	p = parser.NewParser(parser.NewLexer(`if 1 {
	5
} else {
}
1`))
	res = p.Parse()

	compiler.Compile(res)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 1,
		grammar.JumpIfFalse, 1, 15,
		grammar.Advance,
		grammar.Integer, 1, 5, grammar.Noop,
		grammar.Back,
		grammar.Jump, 1, 16,
		grammar.Null,
		grammar.Pop,
		grammar.Advance, grammar.Advance, grammar.Advance, grammar.Advance,
		grammar.Integer, 1, 1, grammar.Pop,
	}, compiler.Result)
}

func TestLambda(t *testing.T) {
	compiler := NewCompiler()

	p := parser.NewParser(parser.NewLexer(`|a, b| {
	a + b
}`))
	res := p.Parse()

	compiler.Compile(res)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Lambda,
		1, 2, 1, 'a', 1, 'b',

		1, 15,

		grammar.Magic,
		grammar.Advance,
		grammar.GetVar, 1, 'b',
		grammar.GetVar, 1, 'a',
		grammar.GetInstance, 1, '+',
		grammar.Call, 1, 1,
		grammar.Pop,

		grammar.Advance,
		grammar.Pop,
	}, compiler.Result)

	compiler = NewCompiler()

	p = parser.NewParser(parser.NewLexer(`{
	a + b
}`))
	res = p.Parse()

	compiler.Compile(res)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Closure,
		1, 15,

		grammar.Magic,
		grammar.Advance,
		grammar.GetVar, 1, 'b',
		grammar.GetVar, 1, 'a',
		grammar.GetInstance, 1, '+',
		grammar.Call, 1, 1,
		grammar.Pop,

		grammar.Advance,
		grammar.Pop,
	}, compiler.Result)
}

func TestGetInstance(t *testing.T) {
	compiler := NewCompiler()

	p := parser.NewParser(parser.NewLexer(`1.nonz`))
	res := p.Parse()

	compiler.Compile(res)
	assert.Equal(t, []byte{grammar.Magic,
		grammar.Integer, 1, 1,
		grammar.GetInstance, 4, 'n', 'o', 'n', 'z',
		grammar.Pop,
	}, compiler.Result)
}
