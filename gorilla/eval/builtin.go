package eval

import (
	"../object"
)

var Builtins map[string]object.Object

func init() {
	Builtins = map[string]object.Object{}
	for _, def := range object.Builtins {
		Builtins[def.Name] = def.Builtin
	}
}
