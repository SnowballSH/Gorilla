package objects

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
	assert.Equal(t, 1, len(obj.instanceVariables().names()))
	obj.setInstanceVariables(newEnvironment())
}
