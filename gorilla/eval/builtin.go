package eval

import (
	"../ast"
	"../object"
	"fmt"
	"strconv"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len": {Fn: LenFunc},

	"display": {
		Fn: func(self object.Object, line int, args ...object.Object) object.Object {
			for _, arg := range args {
				_, _ = fmt.Fprintf(OUT, arg.Inspect()+"\n")
			}
			return NULL
		},
	},
}

func LenFunc(_ object.Object, line int, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("[Line %d] Argument mismatch (expected %d, got %d)", line,
			1, len(args))
	}

	return CallAttr(args[0], "_len", line)
}

func NewFunction(
	params []*ast.Identifier,
	body *ast.BlockStatement,
	env *object.Environment,
	line int,
) *object.Function {
	return &object.Function{
		Parameters: params,
		Body:       body,
		Env:        env,
		SLine:      line,
	}
}

func NewInt(value int64, line int) *object.Integer {
	return &object.Integer{
		Value: value,
		SLine: line,
		Attrs: map[string]object.Object{
			"toStr": &object.Builtin{
				Fn: func(self object.Object, line int, args ...object.Object) object.Object {
					return NewString(strconv.Itoa(int(self.(*object.Integer).Value)), line)
				},
			},
		},
	}
}

func NewString(value string, line int) *object.String {
	return &object.String{
		Value: value,
		SLine: line,
		Attrs: map[string]object.Object{
			"strip": &object.Builtin{
				Fn: func(self object.Object, line int, args ...object.Object) object.Object {
					return NewString(strings.TrimSpace(self.(*object.String).Value), line)
				},
			},
			"_len": &object.Builtin{
				Fn: func(self object.Object, line int, args ...object.Object) object.Object {
					return NewInt(int64(len(self.(*object.String).Value)), line)
				},
			},
			"toInt": &object.Builtin{
				Fn: func(self object.Object, line int, args ...object.Object) object.Object {
					k := self.(*object.String).Value
					val, err := strconv.Atoi(k)
					if err != nil {
						return newError("[Line %d] Cannot parse to Int: '%s'", line, k)
					}
					return NewInt(int64(val), line)
				},
			},
		},
	}
}
