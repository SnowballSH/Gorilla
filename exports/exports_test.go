package exports

import (
	"github.com/SnowballSH/Gorilla/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOverall(t *testing.T) {
	x, e := CompileGorilla("123")
	assert.Nil(t, e)

	_, b, e := ExecuteGorillaBytecodeFromSource(x, "123")
	assert.Nil(t, e)

	assert.Equal(t, "123", b.ToString())

	_, b, e = ExecuteGorillaBytecode(x)
	assert.Nil(t, e)

	assert.Equal(t, "123", b.ToString())

	_, b, e = ExecuteGorillaBytecodeFromSourceAndEnv(x, "123", runtime.NewEnvironment())
	assert.Nil(t, e)

	assert.Equal(t, "123", b.ToString())

	x, e = CompileGorilla(")")
	assert.NotNil(t, e)

	x, e = CompileGorilla("a")
	assert.Nil(t, e)

	_, _, e = ExecuteGorillaBytecodeFromSource(x, "a")
	assert.NotNil(t, e)

	_, _, e = ExecuteGorillaBytecode(x)
	assert.NotNil(t, e)

	_, b, e = ExecuteGorillaBytecodeFromSourceAndEnv(x, "a", runtime.NewEnvironmentWithStore(
		map[string]runtime.BaseObject{"a": runtime.NewInteger(0)},
	))
	assert.Nil(t, e)

	assert.Equal(t, "0", b.ToString())
}
