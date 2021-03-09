package objects

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv(t *testing.T) {
	integer := NewInteger(56789)
	env := newEnvironment()
	env.set("x", integer)
	assert.Equal(t, 1, len(env.store))
	v, o := env.get("x")
	assert.True(t, o)
	assert.EqualValues(t, 56789, v.Value())
	assert.Equal(t, 1, len(env.names()))
	assert.Equal(t, 1, len(env.copy().copy().names()))
}
