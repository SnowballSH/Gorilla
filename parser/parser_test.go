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
}

func TestSimpleParse(t *testing.T) {
	p := NewParser(NewLexer("\n65500\n"))
	res := p.Parse()
	assert.Equal(t, 1, len(res))
	var n *string = nil
	assert.Equal(t, n, p.error)
	assert.Equal(t, "65500;", res[0].String())

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
}
