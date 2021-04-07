package runtime

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool(t *testing.T) {
	assert.Equal(t, "true", GorillaTrue.ToString())
	assert.Equal(t, "true", GorillaTrue.Inspect())
	assert.True(t, GorillaTrue.IsTruthy())

	assert.True(t, GorillaTrue.EqualTo(GorillaTrue))
	assert.False(t, GorillaTrue.EqualTo(BoolClass))

	assert.Equal(t, "false", GorillaFalse.ToString())
	assert.Equal(t, "false", GorillaFalse.Inspect())
	assert.False(t, GorillaFalse.IsTruthy())

	assert.True(t, GorillaFalse.EqualTo(GorillaFalse))
	assert.False(t, GorillaFalse.EqualTo(BoolClass))
}

func TestBoolFuncs(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.GetVar,
		4, 't', 'r', 'u', 'e',
		grammar.GetInstance,
		2, '!', '@',
		grammar.Call,
		1, 0x00,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "false")

	vm = NewVM([]byte{grammar.Magic,
		grammar.GetVar,
		5, 'f', 'a', 'l', 's', 'e',
		grammar.GetInstance,
		2, '!', '@',
		grammar.Call,
		1, 0x00,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "true")
}
