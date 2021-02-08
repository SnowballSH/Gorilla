package stdlib

import (
	"Gorilla/object"
	"math/rand"
)

var RandomNamespace = object.NewNameSpace("random", map[string]object.BaseObject{
	"randint": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			v := args[0].Value().(int)
			if v < 0 {
				return object.NewError("Function 'randint' expects a positive integer, not negative", line)
			}
			return object.NewInteger(rand.Intn(v), line)
		},
		[][]string{{object.INTEGER}},
	),
})
