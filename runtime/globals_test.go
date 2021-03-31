package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobals(t *testing.T) {
	o, err := Global.Store["print"].Call(nil, NewString("abc"), NewInteger(1))
	assert.Nil(t, err)
	assert.Equal(t, o.ToString(), "abc 1")

	o, err = Global.Store["str"].Call(nil, NewInteger(1))
	assert.Nil(t, err)
	assert.Equal(t, o.ToString(), "1")

	o, err = Global.Store["str"].Call(nil)
	assert.NotNil(t, err)
}
