package runtime

import "fmt"

// Every Gorilla Object and Class inherits from this
type BaseObject interface {
	Class() *RClass
	Value() interface{}
	ToString() string
	Inspect() string
	InstanceVariableGet(string) (BaseObject, bool)
	InstanceVariables() *Environment

	IsTruthy() bool
	EqualTo(BaseObject) bool

	Call(self BaseObject, args ...BaseObject) (BaseObject, error)

	Parent() BaseObject
	SetParent(object BaseObject)
}

// Object struct holds a normal object
type Object struct {
	RClass        *RClass
	InternalValue interface{}
	ToStringFunc  func(self *Object) string
	InspectFunc   func(self *Object) string
	IsTruthyFunc  func(self *Object) bool
	EqualToFunc   func(self *Object, other BaseObject) bool
	CallFunc      CallFuncType
	ParentObj     BaseObject
}

func (o *Object) Class() *RClass {
	return o.RClass
}

func (o *Object) Value() interface{} {
	return o.InternalValue
}

func (o *Object) ToString() string {
	return o.ToStringFunc(o)
}

func (o *Object) Inspect() string {
	return o.InspectFunc(o)
}

func (o *Object) InstanceVariableGet(s string) (BaseObject, bool) {
	return o.Class().GetInstance(s)
}

func (o *Object) InstanceVariables() *Environment {
	return o.Class().InstanceVars
}

func (o *Object) IsTruthy() bool {
	return o.IsTruthyFunc(o)
}

func (o *Object) EqualTo(object BaseObject) bool {
	return o.EqualToFunc(o, object)
}

func (o *Object) Call(self BaseObject, args ...BaseObject) (BaseObject, error) {
	return o.CallFunc(self, args...)
}

func (o *Object) Parent() BaseObject {
	return o.ParentObj
}

func (o *Object) SetParent(parent BaseObject) {
	o.ParentObj = parent
}

// Helper

type CallFuncType func(self BaseObject, args ...BaseObject) (BaseObject, error)

type ConvertFuncType func(x BaseObject) (*Object, error)

var NotCallable = func(self BaseObject, args ...BaseObject) (BaseObject, error) {
	return nil, fmt.Errorf(fmt.Sprintf("'%s' is not callable", self.ToString()))
}
