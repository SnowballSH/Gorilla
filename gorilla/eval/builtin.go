package eval

import (
	"../object"
	"fmt"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(line int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("[Line %d] Argument mismatch (expected %d, got %d)", line,
					1, len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},

	"display": {
		Fn: func(line int, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
