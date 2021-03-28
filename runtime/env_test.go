package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv(t *testing.T) {
	integer := NewInteger(56789)
	env := NewEnvironment()
	env.Set("x", integer)
	assert.Equal(t, 1, len(env.Store))
	v, o := env.Get("x")
	assert.True(t, o)
	assert.EqualValues(t, 56789, v.Value())
	assert.Equal(t, 1, len(env.Names()))
	assert.Equal(t, 1, len(env.Copy().Copy().Names()))
}
