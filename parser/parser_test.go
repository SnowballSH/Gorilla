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
	p := NewParser(NewLexer("\n65500;\n"))
	res := p.Parse()
	assert.Equal(t, 1, len(res))
	assert.Nil(t, p.Error)
	assert.Equal(t, "65500;", res[0].String())

	p = NewParser(NewLexer("'abcde'"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))
	assert.Nil(t, p.Error)
	assert.Equal(t, "'abcde';", res[0].String())
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

	p = NewParser(NewLexer("3 * 9 + 4 - 5 % 6 * 999"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "(((3 * 9) + 4) - ((5 % 6) * 999));", res[0].String())

	p = NewParser(NewLexer("3 * (9 + 4) - 5 * (6 * 999)"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "((3 * (9 + 4)) - (5 * (6 * 999)));", res[0].String())

	p = NewParser(NewLexer("+-2"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "(+(-2));", res[0].String())
}

func TestComparison(t *testing.T) {
	p := NewParser(NewLexer("2 == 3\n3 != 4"))
	res := p.Parse()
	assert.Equal(t, 2, len(res))

	assert.Equal(t, "(2 == 3);", res[0].String())
	assert.Equal(t, "(3 != 4);", res[1].String())

	p = NewParser(NewLexer("!2"))
	res = p.Parse()
	assert.Equal(t, 1, len(res))

	assert.Equal(t, "(!2);", res[0].String())
}

func TestVar(t *testing.T) {
	p := NewParser(NewLexer("$abc * 3"))
	res := p.Parse()
	assert.Equal(t, "($abc * 3);", res[0].String())

	p = NewParser(NewLexer("a = 1"))
	res = p.Parse()
	assert.Equal(t, "(a = 1);", res[0].String())

	p = NewParser(NewLexer("2 + a = (2 + 3) * 7"))
	res = p.Parse()
	assert.Equal(t, "(2 + (a = ((2 + 3) * 7)));", res[0].String())
}

func TestBlock(t *testing.T) {
	p := NewParser(NewLexer(`if 0 {`))
	p.Parse()
	assert.NotNil(t, p.Error)

	p = NewParser(NewLexer(`if 0 }`))
	p.Parse()
	assert.NotNil(t, p.Error)

	p = NewParser(NewLexer(`if 0 {0 0}`))
	p.Parse()
	assert.NotNil(t, p.Error)

	p = NewParser(NewLexer(`if 0 {0;;;;}`))
	p.Parse()
	assert.Nil(t, p.Error)
}

func TestIfElse(t *testing.T) {
	p := NewParser(NewLexer(`
if 1 + 2 {
	print("Hello");
	3 + 1
}
`))
	res := p.Parse()
	assert.Equal(t, `(if (1 + 2) {
(print('Hello'));
(3 + 1);
} else {
});`, res[0].String())

	p = NewParser(NewLexer(`
if 1 + 2 {
	print("Hello");
	3 + 1
} else {
	k
}
`))
	res = p.Parse()
	assert.Equal(t, `(if (1 + 2) {
(print('Hello'));
(3 + 1);
} else {
k;
});`, res[0].String())
}

func TestLambda(t *testing.T) {
	p := NewParser(NewLexer(`
|a, b|
	a + b
`))
	res := p.Parse()
	assert.Equal(t, `(|a, b| {
(a + b);
});`, res[0].String())

	p = NewParser(NewLexer(`|| {}`))
	res = p.Parse()
	assert.Equal(t, `(|| {
});`, res[0].String())

	p = NewParser(NewLexer(`|1| {}`))
	p.Parse()
	assert.NotNil(t, p.Error)

	p = NewParser(NewLexer(`|a b| {}`))
	p.Parse()
	assert.NotNil(t, p.Error)

	p = NewParser(NewLexer(`|`))
	p.Parse()
	assert.NotNil(t, p.Error)
}

func TestCall(t *testing.T) {
	p := NewParser(NewLexer("$abc(1,\n 2, 3\n)"))
	res := p.Parse()
	assert.Equal(t, "($abc(1, 2, 3));", res[0].String())

	p = NewParser(NewLexer("$abc(1, 2, 3,\n)"))
	res = p.Parse()
	assert.Equal(t, "($abc(1, 2, 3));", res[0].String())

	p = NewParser(NewLexer("$abc(1 2)"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

$abc(1 2)
       ^
Expected ',', got '2'
`), *p.Error)

	p = NewParser(NewLexer("$abc(1"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

$abc(1
      ^
Expected ',', got End of File
`), *p.Error)

	p = NewParser(NewLexer("$abc("))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

$abc(
     ^
Expected ')', got End of File
`), *p.Error)
}

func TestError(t *testing.T) {
	var p *Parser

	p = NewParser(NewLexer("65500 123"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

65500 123
      ^
Expected newline or ;, got '123'
`), *p.Error)

	p = NewParser(NewLexer("99999999999999999999999"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

99999999999999999999999
^
Could not parse '99999999999999999999999' as 64-bit integer
`), *p.Error)

	p = NewParser(NewLexer("✔"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 1:

✔
^
Unexpected '✔'
`), *p.Error)

	p = NewParser(NewLexer("(\n1"))
	p.Parse()
	assert.Equal(t, strings.TrimSpace(`
Error in line 2:

1
 ^
Expected ')', got End of File
`), *p.Error)
}
