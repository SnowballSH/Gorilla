package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNull(t *testing.T) {
	assert.Equal(t, "null", Null.ToString())
	assert.Equal(t, "null", Null.Inspect())
	assert.False(t, Null.IsTruthy())

	assert.True(t, Null.EqualTo(Null))
	assert.False(t, Null.EqualTo(NullClass))

	a, _ := NullClass.NewFunc(nil, Null)
	assert.Nil(t, a)
}
