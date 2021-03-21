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
		grammar.Back,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.NotNil(t, vm.LastPopped)
	assert.Equal(t, "624485", vm.LastPopped.ToString())
	assert.Equal(t, 1, vm.line)
}

func TestVMGetAttribute(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Integer, 1, 0x01,
		grammar.GetInstance,
		1, '?',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)
	assert.Equal(t, errors.MakeVMError("Attribute '?' does not exist on '1' (class 'Integer')", 0), vm.Error)
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
}

func TestVMVariables(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.SetVar, 3, 'a', 'b', 'c',
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "3")

	vm = NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.SetVar, 3, 'a', 'b', 'c',
		grammar.Pop,
		grammar.GetVar, 3, 'a', 'b', 'c',
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "3")

	vm = NewVMWithStore([]byte{grammar.Magic,
		grammar.GetVar, 3, 'a', 'b', 'd', // undefined
		grammar.Pop,
	}, vm.Environment)
	vm.Run()

	assert.NotNil(t, vm.Error)
	assert.Equal(t, errors.MakeVMError("Variable 'abd' is not defined", 0), vm.Error)
}

func TestCall(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.Integer, 1, 0x03,
		grammar.Call,
		1, 0x00,
		grammar.Pop,
	})
	vm.Run()

	assert.NotNil(t, vm.Error)
	assert.Equal(t, errors.MakeVMError("'3' is not callable", 0), vm.Error)
}

func TestVMError(t *testing.T) {
	vm := NewVM([]byte{})
	vm.Run()

	assert.NotNil(t, vm.Error)
	assert.Equal(t, errors.MakeVMError("Not a valid Gorilla bytecode", 0), vm.Error)
}
