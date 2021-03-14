package runtime

// Every Gorilla Object and Class inherits from this
type BaseObject interface {
	Class() *RClass
	Value() interface{}
	ToString() string
	Inspect() string
	InstanceVariableGet(string) (BaseObject, bool)
	InstanceVariableSet(string, BaseObject) BaseObject
	InstanceVariables() *environment
	SetInstanceVariables(*environment)

	IsTruthy() bool
	EqualTo(BaseObject) bool
}

// Object struct holds a normal object
type Object struct {
	RClass        *RClass
	Instances     *environment
	InternalValue interface{}
	ToStringFunc  func(self *Object) string
	InspectFunc   func(self *Object) string
	IsTruthyFunc  func(self *Object) bool
	EqualToFunc   func(self *Object, other BaseObject) bool
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
	return o.Instances.get(s)
}

func (o *Object) InstanceVariableSet(s string, object BaseObject) BaseObject {
	return o.Instances.set(s, object)
}

func (o *Object) InstanceVariables() *environment {
	return o.Instances
}

func (o *Object) SetInstanceVariables(e *environment) {
	o.Instances = e
}

func (o *Object) IsTruthy() bool {
	return o.IsTruthyFunc(o)
}

func (o *Object) EqualTo(object BaseObject) bool {
	return o.EqualToFunc(o, object)
}