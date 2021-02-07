package object

import (
	"Gorilla/code"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

const (
	ERROR = "Error"

	ANY = "Any"

	BUILTINFUNCTION = "Builtin Function"
	FUNCTION        = "Function"

	INTEGER = "Integer"
	FLOAT   = "Float"
	BOOLEAN = "Boolean"
	STRING  = "String"
	ARRAY   = "Array"
	HASH    = "Hash"
	NULL    = "Null"
)

// Every Gorilla Object implements this
type BaseObject interface {
	Type() string
	Inspect() string
	Debug() string
	Line() int
	Value() interface{}
	FindMethod(name string) (BaseObject, BaseObject)
	SetMethod(name string, value BaseObject)
	Call(env *Environment, self *Object, params []BaseObject, line int) BaseObject
	Parent() BaseObject
	SetParent(obj *Object) BaseObject
}

// Implements BaseObject
type Object struct {
	TT            string
	InternalValue interface{}
	InspectValue  func(self BaseObject) string
	DebugValue    func(self BaseObject) string
	SLine         int
	Methods       map[string]BaseObject
	CallFunc      func(env *Environment, self *Object, args []BaseObject, line int) BaseObject
	ParentObj     BaseObject
}

func (o *Object) Type() string {
	return o.TT
}

func (o *Object) Inspect() string {
	return o.InspectValue(o)
}

func (o *Object) Debug() string {
	return o.DebugValue(o)
}

func (o *Object) Line() int {
	return o.SLine
}

func (o *Object) Value() interface{} {
	return o.InternalValue
}

func (o *Object) FindMethod(name string) (BaseObject, BaseObject) {
	v, ok := o.Methods[name]
	if !ok {
		return nil, NewError(fmt.Sprintf("Method not found: %s on type '%s'", name, o.Type()), o.Line())
	}
	return CopyObject(v.(*Object)), nil
}

func (o *Object) SetMethod(name string, value BaseObject) {
	o.Methods[name] = value
}

func (o *Object) Call(env *Environment, self *Object, args []BaseObject, line int) BaseObject {
	return o.CallFunc(env, self, args, line)
}

func (o *Object) Parent() BaseObject {
	return o.ParentObj
}

func (o *Object) SetParent(obj *Object) BaseObject {
	o.ParentObj = obj
	return o
}

// Helper function, creates a new Object
func NewObject(
	TT string,
	InternalValue interface{},
	InspectValue func(self BaseObject) string,
	DebugValue func(self BaseObject) string,
	SLine int,
	Methods map[string]BaseObject,
	CallFunc func(env *Environment, self *Object, args []BaseObject, line int) BaseObject,
	Parent BaseObject,
) *Object {
	if CallFunc == nil {
		CallFunc = func(env *Environment, self *Object, args []BaseObject, line int) BaseObject {
			return NewError(
				fmt.Sprintf("Type '%s' is not Callable", TT),
				SLine,
			)
		}
	}

	var mts = map[string]BaseObject{}

	for n, v := range BaseObjectBuiltins {
		mts[n] = v
	}

	for n, v := range Methods {
		mts[n] = v
	}

	return &Object{
		TT:            TT,
		InternalValue: InternalValue,
		InspectValue:  InspectValue,
		DebugValue:    DebugValue,
		SLine:         SLine,
		Methods:       mts,
		CallFunc:      CallFunc,
		ParentObj:     Parent,
	}
}

// Base ERROR Type
func NewError(
	value string,
	line int,
) *Object {
	return NewObject(
		ERROR,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("[Line %d] %s", self.Line()+1, self.Value().(string))
		},
		func(self BaseObject) string {
			return fmt.Sprintf("Gorilla Error: [Line %d] %s", self.Line()+1, self.Value().(string))
		},
		line,
		map[string]BaseObject{},
		nil,
		nil,
	)
}

// Base BUILTINFUNCTION Type
func NewBuiltinFunction(
	value func(self *Object, env *Environment, args []BaseObject, line int) BaseObject,
	params [][]string,
) *Object {
	return NewObject(
		BUILTINFUNCTION,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("Builtin Function")
		},
		func(self BaseObject) string {
			return fmt.Sprintf("Builtin Function [%p]", self.(*Object))
		},
		0,
		map[string]BaseObject{},
		func(env *Environment, self *Object, args []BaseObject, lline int) BaseObject {
			if params != nil {
				// Argument
				if len(args) != len(params) {
					return NewError(
						fmt.Sprintf("Argument amount mismatch: Expected %d, got %d", len(params), len(args)),
						lline,
					)
				}

				// Type Checking
				for i, v := range args {
					ok := false
					for _, vv := range params[i] {
						if vv == ANY {
							ok = true
							break
						}
						if v.Type() == vv {
							ok = true
							break
						}
					}
					if ok {
						continue
					}
					return NewError(
						fmt.Sprintf(
							"Argument #%d expected to be one of %s, got Type '%s'",
							i, params[i], v.Type(),
						),
						lline,
					)
				}
			}

			return value(self, env, args, lline)
		},
		nil,
	)
}

// Base INTEGER Type
func NewInteger(
	value int,
	line int,
) *Object {
	return NewObject(
		INTEGER,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("%d", self.Value().(int))
		},
		func(self BaseObject) string {
			return fmt.Sprintf("%d", self.Value().(int))
		},
		line,
		IntegerBuiltins,
		nil,
		nil,
	)
}

// Base INTEGER Type
func NewFloat(
	value float64,
	line int,
) *Object {
	return NewObject(
		FLOAT,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("%f", self.Value().(float64))
		},
		func(self BaseObject) string {
			return fmt.Sprintf("%f", self.Value().(float64))
		},
		line,
		FloatBuiltins,
		nil,
		nil,
	)
}

// Base BOOLEAN Type
func NewBool(
	value bool,
	line int,
) *Object {
	return NewObject(
		BOOLEAN,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("%s", strconv.FormatBool(self.Value().(bool)))
		},
		func(self BaseObject) string {
			return fmt.Sprintf("%s", strconv.FormatBool(self.Value().(bool)))
		},
		line,
		BooleanBuiltins,
		nil,
		nil,
	)
}

// Base STRING Type
func NewString(
	value string,
	line int,
) *Object {
	return NewObject(
		STRING,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("%s", self.Value().(string))
		},
		func(self BaseObject) string {
			return fmt.Sprintf("\"%s\"", self.Value().(string))
		},
		line,
		StringBuiltins,
		nil,
		nil,
	)
}

// Base NULL Type
func NewNull(
	line int,
) *Object {
	return NewObject(
		NULL,
		nil,
		func(self BaseObject) string {
			return "null"
		},
		func(self BaseObject) string {
			return "null"
		},
		line,
		map[string]BaseObject{},
		nil,
		nil,
	)
}

type FunctionValue struct {
	Constants []BaseObject
	Bytecodes []code.Opcode
	Messages  []Message
	Params    []string
	FreeEnv   *Environment
}

// Base FUNCTION Type
func NewFunction(
	value *FunctionValue,
	line int,
) *Object {
	o := NewObject(
		FUNCTION,
		value,
		func(self BaseObject) string {
			return "Function"
		},
		func(self BaseObject) string {
			return fmt.Sprintf("Function [%p]", self)
		},
		line,
		map[string]BaseObject{},
		nil, // in vm
		nil,
	)
	return o
}

func af(self BaseObject) string {
	var k []string
	for _, v := range self.Value().([]BaseObject) {
		k = append(k, v.Debug())
	}
	return "[" + strings.Join(k, ", ") + "]"
}

// Base ARRAY Type
func NewArray(
	value []BaseObject,
	line int,
) *Object {
	return NewObject(
		ARRAY,
		value,
		af,
		af,
		line,
		ArrayBuiltins,
		nil,
		nil,
	)
}

func hf(self BaseObject) string {
	var pairs []string
	for _, value := range self.Value().(map[HashKey]*HashValue) {
		pairs = append(pairs, value.Key.Debug()+": "+value.Value.Debug())
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

type HashKey struct {
	Type      string
	HashedKey uint64
}

type HashValue struct {
	Key   BaseObject
	Value BaseObject
}

// Base HASH Type
func NewHash(
	value map[HashKey]*HashValue,
	line int,
) *Object {
	return NewObject(
		HASH,
		value,
		hf,
		hf,
		line,
		HashBuiltins,
		nil,
		nil,
	)
}

func CopyObject(obj *Object) *Object {
	return NewObject(obj.TT, obj.InternalValue, obj.InspectValue, obj.DebugValue, obj.SLine, obj.Methods, obj.CallFunc, obj.ParentObj)
}

func HashObject(obj BaseObject) (key HashKey, ok bool) {
	switch obj.Type() {
	case INTEGER:
		return HashKey{Type: obj.Type(), HashedKey: uint64(obj.Value().(int))}, true
	case STRING:
		h := fnv.New64a()
		_, _ = h.Write([]byte(obj.Value().(string)))
		return HashKey{Type: obj.Type(), HashedKey: h.Sum64()}, true
	case BOOLEAN:
		var val uint64
		if obj.Value().(bool) {
			val = 1
		} else {
			val = 0
		}

		return HashKey{Type: obj.Type(), HashedKey: val}, true
	default:
		return HashKey{}, false
	}
}
