package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClass(t *testing.T) {
	xclass := &IntegerClass
	class := *xclass
	_, o := class.InstanceVariableGet("x")
	assert.False(t, o)
	assert.True(t, class.EqualTo(IntegerClass))
	assert.False(t, class.EqualTo(NumericClass))
	assert.True(t, class.IsTruthy())

	assert.Equal(t, class, class.Class())
	assert.Equal(t, "Class 'Integer'", class.Value())
	assert.Equal(t, "Class 'Integer'", class.Inspect())

	assert.Equal(t, class.superClass, class.Parent())
	class.SetParent(IntegerClass)

	_, err := class.Call(nil, class)
	assert.NotNil(t, err)
}
