package object

import (
	"../ast"
	"../config"
	"bufio"
	"fmt"
	"io"
	"os"
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
var ArrayAttrs map[string]Object
var BoolAttrs map[string]Object

var TRUE *Boolean
var FALSE *Boolean
var NULL *Null

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
			"print",
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
			"input",
			&Builtin{
				Fn: func(self Object, line int, args ...Object) Object {
					buffer := bufio.NewReader(os.Stdin)

					lineC, _, err := buffer.ReadLine()
					if err != nil && err != io.EOF {
						return NewError("[Line %d] EOF When getting input", line+1)
					}
					return NewString(string(lineC), line)
				},
			},
		},

		{
			"len",
			&Builtin{Fn: func(_ Object, line int, args ...Object) Object {
				if len(args) != 1 {
					return NewError("[Line %d] Argument mismatch (expected %d, got %d)", line+1,
						1, len(args))
				}

				switch arg := args[0].(type) {
				case *String:
					return NewInt(int64(utf8.RuneCountInString(arg.Value)), arg.Line())
				case *Array:
					return NewInt(int64(len(arg.Value)), arg.Line())
				default:
					return NewError("[Line %d] Cannot get length of type '%s'", line+1, arg.Type())
				}
			},
			},
		},
	}

	IntAttrs = map[string]Object{
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				// return NewString(strconv.Itoa(int(self.(*Integer).Value)), line+1)
				return NewString(self.(*Integer).Inspect(), line+1)
			},
		}}

	StrAttrs = map[string]Object{
		"strip": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(strings.TrimSpace(self.(*String).Value), line+1)
			},
		},
		"toInt": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				k := self.(*String).Value
				val, err := strconv.Atoi(k)
				if err != nil {
					return NewError("[Line %d] Cannot parse to Int: '%s'", line+1, k)
				}
				return NewInt(int64(val), line)
			},
		},
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(self.(*String).Inspect(), line+1)
			},
		},
	}

	ArrayAttrs = map[string]Object{
		"push": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				arr := self.(*Array)
				return arr.PushAll(args)
			},
		},
		"pop": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				if len(self.(*Array).Value) < 1 {
					return NewError("[Line %d] Cannot pop empty array", line+1)
				}
				return self.(*Array).PopLast()
			},
		},
		"shift": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				if len(self.(*Array).Value) < 1 {
					return NewError("[Line %d] Cannot shift empty array", line+1)
				}
				return self.(*Array).PopFirst()
			},
		},
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(self.(*Array).Inspect(), line+1)
			},
		},
	}

	BoolAttrs = map[string]Object{
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(self.(*Boolean).Inspect(), line+1)
			},
		},
	}

	TRUE = NewBool(true, 0)

	FALSE = NewBool(false, 0)

	NULL = &Null{}
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

func NewArray(value []Object, line int) *Array {
	return &Array{
		Value: value,
		SLine: line,
		Attrs: ArrayAttrs,
	}
}

func NewBool(value bool, line int) *Boolean {
	return &Boolean{
		Value: value,
		SLine: line,
		Attrs: BoolAttrs,
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
