package helper

import (
	"unicode"
)

var combining = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0300, 0x036f, 1}, // combining diacritical marks
		{0x1ab0, 0x1aff, 1}, // combining diacritical marks extended
		{0x1dc0, 0x1dff, 1}, // combining diacritical marks supplement
		{0x20d0, 0x20ff, 1}, // combining diacritical marks for symbols
		{0xfe20, 0xfe2f, 1}, // combining half marks
	},
}

func ReverseString(s string) string {
	sv := []rune(s)
	rv := make([]rune, 0, len(sv))
	cv := make([]rune, 0)
	for ix := len(sv) - 1; ix >= 0; ix-- {
		r := sv[ix]
		if unicode.In(r, combining) {
			cv = append(cv, r)
		} else {
			rv = append(rv, r)
			rv = append(rv, cv...)
			cv = make([]rune, 0)
		}
	}
	return string(rv)
}

func ReverseBOA(val []interface{}) []interface{} {
	var a []interface{}
	copy(a, val)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return val
}
