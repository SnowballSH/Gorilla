package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClass(t *testing.T) {
	xclass := &IntegerClass
	class := *xclass
	class.InstanceVariableSet("x", NewInteger(6))
	assert.Equal(t, 1, len(class.InstanceVariables().names()))
	x, o := class.InstanceVariableGet("x")
	assert.True(t, o)
	assert.EqualValues(t, 6, x.Value())
	assert.True(t, class.EqualTo(IntegerClass))
	assert.False(t, class.EqualTo(NumericClass))
	assert.True(t, class.IsTruthy())

	class.SetInstanceVariables(newEnvironment())
	assert.Equal(t, class, class.Class())
	assert.Equal(t, "Class 'Integer'", class.Value())
	assert.Equal(t, "Class 'Integer'", class.Inspect())
}
