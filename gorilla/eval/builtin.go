package eval

import (
	"../object"
)

var Builtins map[string]*object.Builtin

func init() {
	Builtins = map[string]*object.Builtin{
		"len":     object.LookupBuiltin("len"),
		"display": object.LookupBuiltin("display"),
	}
}
