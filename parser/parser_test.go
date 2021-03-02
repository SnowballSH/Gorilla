package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser(NewLexer("123 abc"))
	assert.Equal(t, "123", p.cur.Literal)
	assert.Equal(t, "abc", p.peek.Literal)

	p = NewParser(NewLexer(""))
	assert.Equal(t, rune(0), p.l.ch)
	assert.Equal(t, 0, p.l.inputLength)
}

func TestSimpleParse(t *testing.T) {
	p := NewParser(NewLexer("\n65500\n"))
	res := p.Parse()
	assert.Equal(t, 1, len(res))
	var n *string = nil
	assert.Equal(t, n, p.error)
	assert.Equal(t, "65500;", res[0].String())
}

func TestParseBinOp(t *testing.T) {
	p := NewParser(NewLexer("2 + 3\n3 * 4"))
	res := p.Parse()
	assert.Equal(t, 2, len(res))

	assert.Equal(t, "(2 + 3);", res[0].String())
	assert.Equal(t, "(3 * 4);", res[1].String())

	p = NewParser(NewLexer("2 + 3 * 4"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "(2 + (3 * 4));", res[0].String())

	p = NewParser(NewLexer("2 / 3 - 4"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "((2 / 3) - 4);", res[0].String())
}

func TestError(t *testing.T) {
	var p *Parser

	p = NewParser(NewLexer("65500 123"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
65500 123
      ^
Expected newline or ;, got 123
`), *p.error)

	p = NewParser(NewLexer("99999999999999999999999"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
99999999999999999999999
^
Could not parse '99999999999999999999999' as 64-bit integer
`), *p.error)

	p = NewParser(NewLexer("✔"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
✔
^
Unexpected '✔'
`), *p.error)
}
