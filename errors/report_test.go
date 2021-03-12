package errors

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMakeError(t *testing.T) {
	code := `
print()
print )
`

	assert.Equal(t, strings.TrimSpace(`
Error in line 2:

print )
      ^
Invalid syntax
`), MakeError(code, "Invalid syntax", 1, 7, 1))
}

func TestPanic(t *testing.T) {
	assert.Panics(t, func() {
		TestERR("")
	})
	assert.Panics(t, func() {
		TestVMERR("")
	})
}

func TestReport(t *testing.T) {
	assert.Equal(t, &VMERROR{
		message: "??",
		line:    1,
	}, MakeVMError("??", 1))
}
