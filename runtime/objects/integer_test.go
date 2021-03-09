package objects

import (
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
	assert.True(t, integer.isTruthy())
	assert.False(t, integer.equalTo(NewInteger(5)))
	assert.False(t, integer.equalTo(nil))

	wot := NewInteger(0)
	wot.InternalValue = "???"
	assert.False(t, integer.equalTo(wot))
}
