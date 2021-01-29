package object

import (
	"../ast"
	"../code"
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
var FunctionClosureAttrs map[string]Object

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
						_, _ = fmt.Fprint(config.OUT, arg.Inspect()+"\n")
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
			&Builtin{
				Fn: func(_ Object, line int, args ...Object) Object {
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

		{
			"typeof",
			&Builtin{
				Fn: func(_ Object, line int, args ...Object) Object {
					if len(args) != 1 {
						return NewError("[Line %d] Argument mismatch (expected %d, got %d)", line+1,
							1, len(args))
					}
					return NewString(string(args[0].Type()), line)
				},
			},
		},

		{
			"isTruthy",
			&Builtin{
				Fn: func(_ Object, line int, args ...Object) Object {
					if len(args) != 1 {
						return NewError("[Line %d] Argument mismatch (expected %d, got %d)", line+1,
							1, len(args))
					}
					return NewBool(IsTruthy(args[0]), line)
				},
			},
		},

		{
			"exit",
			&Builtin{
				Fn: func(_ Object, line int, args ...Object) Object {
					if len(args) != 1 {
						return NewError("[Line %d] Argument mismatch (expected %d, got %d)", line+1,
							1, len(args))
					}
					if _, ok := args[0].(*Integer); !ok {
						return NewError("[Line %d] Expected Integer, got %s", line+1,
							args[0].Type())
					}
					os.Exit(int(args[0].(*Integer).Value))
					return NULL
				},
			},
		},
	}

	IntAttrs = map[string]Object{
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(self.(*Integer).Inspect(), line)
			},
		},
		"range": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				if len(args) < 1 {
					return NewError("[Line %d] Missing required argument 'end'", line)
				}
				if _, ok := args[0].(*Integer); !ok {
					return NewError("[Line %d] Range End is not Integer", line)
				}
				return NewArray(makeRange(int(self.(*Integer).Value), int(args[0].(*Integer).Value), line), line)
			},
		},
		"nonz": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewBool(self.(*Integer).Value != 0, line)
			},
		},
		"positive": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewBool(self.(*Integer).Value > 0, line)
			},
		},
		"negative": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewBool(self.(*Integer).Value < 0, line)
			},
		},
	}

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
		"ord": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				r := []rune(self.(*String).Value)
				rr := make([]Object, len(r))
				for i, v := range r {
					rr[i] = NewInt(int64(v), line)
				}
				return NewArray(rr, line)
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
		"toInt": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				var v int64
				switch self.(*Boolean).Value {
				case true:
					v = 1
				default:
					v = 0
				}
				return NewInt(v, line+1)
			},
		},
	}

	FunctionClosureAttrs = map[string]Object{
		"toStr": &Builtin{
			Fn: func(self Object, line int, args ...Object) Object {
				return NewString(self.Inspect(), line+1)
			},
		},
	}

	TRUE = NewBool(true, 0)

	FALSE = NewBool(false, 0)

	NULL = NewNull()
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func NewNull() *Null {
	return &Null{Attrs: CopyMap(FunctionClosureAttrs)}
}

func NewInt(value int64, line int) *Integer {
	return &Integer{
		Value: value,
		SLine: line,
		Attrs: CopyMap(IntAttrs),
	}
}

func NewString(value string, line int) *String {
	return &String{
		Value: value,
		SLine: line,
		Attrs: CopyMap(StrAttrs),
	}
}

func NewArray(value []Object, line int) *Array {
	return &Array{
		Value: value,
		SLine: line,
		Attrs: CopyMap(ArrayAttrs),
	}
}

func NewBool(value bool, line int) *Boolean {
	return &Boolean{
		Value: value,
		SLine: line,
		Attrs: CopyMap(BoolAttrs),
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
		Attrs:      CopyMap(FunctionClosureAttrs),
	}
}

func NewCompiledFunction(
	ins code.Instructions,
	line, locals, params int,
) *CompiledFunction {
	return &CompiledFunction{
		Instructions:  ins,
		NumLocals:     locals,
		NumParameters: params,
		SLine:         line,
		Attrs:         CopyMap(FunctionClosureAttrs),
	}
}

func NewClosure(
	Fn *CompiledFunction,
	Free []Object,
	SLine int,
) *Closure {
	return &Closure{
		Fn:    Fn,
		Free:  Free,
		SLine: SLine,
		Attrs: CopyMap(FunctionClosureAttrs),
	}
}

func makeRange(min, max, line int) []Object {
	if max <= min {
		return []Object{}
	}
	a := make([]Object, max-min+1)
	for i := range a {
		a[i] = NewInt(int64(min+i), line)
	}
	return a
}

func CopyMap(m map[string]Object) map[string]Object {
	cp := make(map[string]Object)
	for k, v := range m {
		cp[k] = v
	}

	return cp
}

func IsTruthy(obj Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
