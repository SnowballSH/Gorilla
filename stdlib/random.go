package stdlib

import (
	"Gorilla/object"
	"math/rand"
)

// Random util
var RandomNamespace = object.NewNameSpace("random", map[string]object.BaseObject{
	"intRange": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			v := args[0].Value().(int)
			if v <= 0 {
				return object.NewError("Function 'intRange' expects a non-negative integer, not negative", line)
			}
			return object.NewInteger(rand.Intn(v), line)
		},
		[][]string{{object.INTEGER}},
	),

	"float": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(rand.Float64(), line)
		},
		[][]string{},
	),

	"int": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewInteger(rand.Int(), line)
		},
		[][]string{},
	),
})
