package errors

import (
	"strings"
)

type PARSINGERROR struct {
}

func MakeError(code, why string, line, char, e int) string {
	return strings.Split(strings.ReplaceAll(strings.TrimSpace(code), "\r", ""), "\n")[line] +
		"\n" + strings.Repeat(" ", char-e) + "^" + "\n" + why
}
