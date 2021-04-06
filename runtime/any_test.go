package runtime

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAnyFunc(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		2, '=', '=',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "0")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		2, '!', '=',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "1")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x03,
		grammar.GetInstance,
		2, '=', '=',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "1")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x03,
		grammar.GetInstance,
		2, '!', '=',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "0")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.GetInstance,
		2, '=', '=',
		grammar.Call,
		1, 0x00,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.GetInstance,
		2, '!', '=',
		grammar.Call,
		1, 0x00,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)
}
