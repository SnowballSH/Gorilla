package object

import (
	"../ast"
	"../config"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

var Builtins []struct {
	Name    string
	Builtin Object
}

var IntAttrs map[string]Object
var StrAttrs map[string]Object
var AllAttrs []struct {
	N string
	T []string
}

func init() {
	Builtins = []struct {
		Name    string
		Builtin Object
	}{
		{
			"null",
			NULL,
		},

		{
			"display",
			&Builtin{
				Fn: func(self Object, line int, args ...Object) Object {
					for _, arg := range args {
						_, _ = fmt.Fprintf(config.OUT, arg.Inspect()+"\n")
					}
					return NULL
				},
			},
		},

		{
			"len",
			&Builtin{Fn: func(_ Object, line int, args ...Object) Object {
				if len(args) != 1 {
					return NewError("[Line %d] Argument mismatch (expected %d, got %d)", line,
						1, len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return NewInt(int64(utf8.RuneCountInString(arg.Value)), arg.Line())
				default:
					return NewError("[Line %d] Cannot get length of type '%s'", line, arg.Type())
				}
			},
			},
		},
	}

	IntAttrs = map[string]Object{
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(strconv.Itoa(int(self.(*Integer).Value)), line)
			},
		}}

	StrAttrs = map[string]Object{
		"strip": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(strings.TrimSpace(self.(*String).Value), line)
			},
		},
		"toInt": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				k := self.(*String).Value
				val, err := strconv.Atoi(k)
				if err != nil {
					return NewError("[Line %d] Cannot parse to Int: '%s'", line, k)
				}
				return NewInt(int64(val), line)
			},
		},
	}

	AllAttrs = []struct {
		N string
		T []string
	}{
		{"toStr", []string{INTEGER}},

		{"toInt", []string{STRING}},

		{"strip", []string{STRING}},
	}
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func NewInt(value int64, line int) *Integer {
	return &Integer{
		Value: value,
		SLine: line,
		Attrs: IntAttrs,
	}
}

func NewString(value string, line int) *String {
	return &String{
		Value: value,
		SLine: line,
		Attrs: StrAttrs,
	}
}

func NewFunction(
	params []*ast.Identifier,
	body *ast.BlockStatement,
	env *Environment,
	line int,
) *Function {
	return &Function{
		Parameters: params,
		Body:       body,
		Env:        env,
		SLine:      line,
	}
}
