package stdlib

import (
	"Gorilla/object"
	"math"
)

/// Math Builtin, maps to Golang's math module for speed
var MathNamespace = object.NewNameSpace("math", map[string]object.BaseObject{
	"mod": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(math.Mod(args[0].Value().(float64), args[1].Value().(float64)), line)
		},
		[][]string{{object.FLOAT}, {object.FLOAT}},
	),
	"pow": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(math.Pow(args[0].Value().(float64), args[1].Value().(float64)), line)
		},
		[][]string{{object.FLOAT}, {object.FLOAT}},
	),

	"sin": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(math.Sin(args[0].Value().(float64)), line)
		},
		[][]string{{object.FLOAT}},
	),
	"cos": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(math.Cos(args[0].Value().(float64)), line)
		},
		[][]string{{object.FLOAT}},
	),
	"tan": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			return object.NewFloat(math.Tan(args[0].Value().(float64)), line)
		},
		[][]string{{object.FLOAT}},
	),
})
