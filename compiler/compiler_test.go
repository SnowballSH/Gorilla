package compiler

import (
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasic(t *testing.T) {
	comp := NewCompiler(parser.NewParser(parser.NewLexer("120")).Parse())
	comp.compile()
	assert.Equal(t, []byte{0x69, 0x00, 0x02, 248, 0}, comp.Result)
}
