package object

import (
	"fmt"
	"strconv"
)

const (
	ERROR = "Error"

	MAIN = "MAIN"
	ANY  = "Any"

	BUILTINFUNCTION = "Builtin Function"

	INTEGER = "Integer"
	BOOLEAN = "Boolean"
)

// Every Gorilla Object implements this
type BaseObject interface {
	Type() string
	Inspect() string
	Line() int
	Value() interface{}
	FindMethod(name string) (BaseObject, BaseObject)
	SetMethod(name string, value BaseObject)
	Call(env *Environment, self *Object, params []BaseObject, line int) BaseObject
	Parent() BaseObject
	SetParent(obj BaseObject) BaseObject
}

// Implements BaseObject
type Object struct {
	TT            string
	InternalValue interface{}
	InspectValue  func(self BaseObject) string
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

func (o *Object) Line() int {
	return o.SLine
}

func (o *Object) Value() interface{} {
	return o.InternalValue
}

func (o *Object) FindMethod(name string) (BaseObject, BaseObject) {
	v, ok := o.Methods[name]
	if !ok {
		return nil, NewError(fmt.Sprintf("Method not found: %s", name), o.Line())
	}
	return v, nil
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

func (o *Object) SetParent(obj BaseObject) BaseObject {
	o.ParentObj = obj
	return o
}

// Helper function, creates a new Object
func NewObject(
	TT string,
	InternalValue interface{},
	InspectValue func(self BaseObject) string,
	SLine int,
	Methods map[string]BaseObject,
	CallFunc func(env *Environment, self *Object, args []BaseObject, line int) BaseObject,
	Parent BaseObject,
) *Object {
	if CallFunc == nil {
		CallFunc = func(env *Environment, self *Object, args []BaseObject, line int) BaseObject {
			return NewError(
				fmt.Sprintf("Type %s is not Callable", self.Type()),
				self.Line(),
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
		line,
		map[string]BaseObject{},
		nil,
		nil,
	)
}

// The MAIN Object
func NewMain() *Object {
	return NewObject(
		MAIN, "",
		func(self BaseObject) string { return "MAIN" }, 0,
		map[string]BaseObject{}, nil,
		nil,
	)
}

// Base BUILTINFUNCTION Type
func NewBuiltinFunction(
	value func(self *Object, args []BaseObject, line int) BaseObject,
	params [][]string,
) *Object {
	return NewObject(
		BUILTINFUNCTION,
		value,
		func(self BaseObject) string {
			return fmt.Sprintf("Builtin Function [%p]", self)
		},
		0,
		map[string]BaseObject{},
		func(env *Environment, self *Object, args []BaseObject, lline int) BaseObject {
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
						"Argument #%d expected to be Type '%s', got Type '%s'",
						i, params[i], v.Type(),
					),
					lline,
				)
			}

			return value(self, args, lline)
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
		line,
		IntegerBuiltins,
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
		line,
		map[string]BaseObject{},
		nil,
		nil,
	)
}
