package runtime

import (
	"github.com/SnowballSH/Gorilla/grammar"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	str := NewString("abc")
	assert.EqualValues(t, "abc", str.InternalValue)
	assert.EqualValues(t, "abc", str.Value())
	assert.Equal(t, "abc", str.ToString())
	assert.Equal(t, "'abc'", str.Inspect())
	assert.Equal(t, StringClass, str.Class())
	assert.True(t, str.IsTruthy())
	assert.False(t, str.EqualTo(NewString("?")))
	assert.False(t, str.EqualTo(nil))

	wot := NewInteger(0)
	assert.False(t, str.EqualTo(wot))
}

func TestStringMethods(t *testing.T) {
	vm := NewVM([]byte{grammar.Magic,
		grammar.String, 1, 'a',
		grammar.String, 1, 'b',
		grammar.GetInstance,
		1, '+',
		grammar.Call,
		1, 0x01,
		grammar.Pop,
	})
	vm.Run()

	assert.Nil(t, vm.Error)
	assert.Equal(t, vm.LastPopped.ToString(), "ba")
}
