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
	obj.InstanceVariableSet("x", NewInteger(8))
	x, o = obj.InstanceVariableGet("x")
	assert.True(t, o)
	assert.NotNil(t, x)
	obj.InstanceVariables().names()
	obj.SetInstanceVariables(newEnvironment())

	assert.Nil(t, obj.Parent())
	obj.SetParent(obj)
	assert.Equal(t, obj, obj.Parent())
}
