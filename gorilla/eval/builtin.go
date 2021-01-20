package eval

import (
	"../object"
)

var Builtins map[string]object.Object

func init() {
	Builtins = map[string]object.Object{
		"len":     object.LookupBuiltin("len"),
		"display": object.LookupBuiltin("display"),
		"null":    object.LookupBuiltin("null"),
	}
}
