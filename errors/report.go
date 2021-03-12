package errors

import (
	"fmt"
	"strings"
)

type PARSINGERROR byte

func MakeError(code, why string, line, char, e int) string {
	return fmt.Sprintf("Error in line %d:\n\n", line+1) + strings.Split(strings.ReplaceAll(strings.TrimSpace(code), "\r", ""), "\n")[line] +
		"\n" + strings.Repeat(" ", char-e) + "^" + "\n" + why
}

type VMERROR struct {
	message string
	line    int
}

func MakeVMError(why string, l int) *VMERROR {
	return &VMERROR{
		message: why,
		line:    l,
	}
}

func TestERR(r interface{}) {
	if _, ok := r.(PARSINGERROR); !ok {
		panic(r)
	}
}

func TestVMERR(r interface{}) {
	if _, ok := r.(*VMERROR); !ok {
		panic(r)
	}
}
