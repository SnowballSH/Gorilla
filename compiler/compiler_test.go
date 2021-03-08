package compiler

import (
	"github.com/SnowballSH/Gorilla/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicNumbers(t *testing.T) {
	comp, ok := Compile(parser.NewParser(parser.NewLexer("120")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{0x69, 0x00, 0x02, 0xf8, 0x00}, comp)

	comp, ok = Compile(parser.NewParser(parser.NewLexer("69420")).Parse())
	assert.True(t, ok)
	assert.Equal(t, []byte{0x69, 0x00, 0x03, 0xac, 0x9e, 0x04}, comp)
}
