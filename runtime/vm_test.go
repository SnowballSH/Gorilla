package runtime

import (
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicVM(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 3, 0xAC, 0x9E, 0x04, grammar.Pop,
		grammar.Advance, grammar.Advance,
		grammar.Integer, 3, 0xE5, 0x8E, 0x26, grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.error)
	assert.NotNil(t, vm.lastPopped)
	assert.Equal(t, "624485", vm.lastPopped.ToString())
}

func TestVMError(t *testing.T) {
	vm := NewVM([]byte{})
	vm.Run()

	assert.NotNil(t, vm.error)
	assert.Equal(t, errors.MakeVMError("Not a valid Gorilla bytecode", 0), vm.error)
}
