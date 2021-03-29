package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObj(t *testing.T) {
	obj := NewInteger(0)
	x, o := obj.InstanceVariableGet("x")
	assert.False(t, o)
	assert.Nil(t, x)
	obj.InstanceVariables().Names()

	assert.Nil(t, obj.Parent())
	obj.SetParent(obj)
	assert.Equal(t, obj, obj.Parent())
}
