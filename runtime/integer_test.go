package runtime

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInteger(t *testing.T) {
	integer := NewInteger(56789)
	assert.EqualValues(t, 56789, integer.InternalValue)
	assert.EqualValues(t, 56789, integer.Value())
	assert.Equal(t, "56789", integer.ToString())
	assert.Equal(t, "56789", integer.Inspect())
	assert.Equal(t, IntegerClass, integer.Class())
	assert.True(t, integer.IsTruthy())
	assert.False(t, integer.EqualTo(NewInteger(5)))
	assert.False(t, integer.EqualTo(nil))

	wot := NewString("")
	assert.False(t, integer.EqualTo(wot))
}

func TestIntegerBinOp(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '*',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "3")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x06,
		grammar.GetInstance,
		1, '/',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "2")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x00,
		grammar.Integer, 1, 0x06,
		grammar.GetInstance,
		1, '/',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x06,
		grammar.GetInstance,
		1, '+',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "9")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x06,
		grammar.GetInstance,
		1, '-',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "3")

	vm = NewVM([]byte{grammar.Magic,
		grammar.String, 1, '1',
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '*',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)

	vm = NewVM([]byte{grammar.Magic,
		grammar.String, 1, '1',
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '/',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)

	vm = NewVM([]byte{grammar.Magic,
		grammar.String, 1, '1',
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '+',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)

	vm = NewVM([]byte{grammar.Magic,
		grammar.String, 1, '1',
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '-',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)
}
