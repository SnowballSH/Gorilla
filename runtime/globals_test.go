package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobals(t *testing.T) {
	o, err := Global.store["print"].Call(nil, NewString("abc"), NewInteger(1))
	assert.Nil(t, err)
	assert.Equal(t, o.ToString(), "abc 1")
}
