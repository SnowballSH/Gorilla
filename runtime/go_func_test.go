package runtime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoFunc(t *testing.T) {
	f := NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
		return NewInteger(666), nil
	})

	assert.Equal(t, "Native Function", f.ToString())
	assert.Equal(t, fmt.Sprintf("Native Function %p", f), f.Inspect())
	assert.True(t, f.IsTruthy())
	assert.True(t, intIns.store["+"].EqualTo(intIns.store["+"]))
	assert.False(t, intIns.store["+"].EqualTo(intIns.store["-"]))
}
