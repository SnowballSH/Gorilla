package runtime

// Defines a Class, or Type
type RClass struct {
	Name       string
	Instances  *Environment
	superClass *RClass

	NewFunc CallFuncType
}

func (R *RClass) Class() *RClass {
	return R
}

func (R *RClass) Value() interface{} {
	return R.ToString()
}

func (R *RClass) ToString() string {
	return "Class '" + R.Name + "'"
}

func (R *RClass) Inspect() string {
	return R.ToString()
}

func (R *RClass) InstanceVariableGet(s string) (BaseObject, bool) {
	return R.Instances.get(s)
}

func (R *RClass) InstanceVariableSet(s string, object BaseObject) BaseObject {
	return R.Instances.set(s, object)
}

func (R *RClass) InstanceVariables() *Environment {
	return R.Instances
}

func (R *RClass) SetInstanceVariables(e *Environment) {
	R.Instances = e
}

func (R *RClass) IsTruthy() bool {
	return true
}

func (R *RClass) EqualTo(object BaseObject) bool {
	return object == R
}

func (R *RClass) Call(a BaseObject, b ...BaseObject) (BaseObject, error) {
	return R.NewFunc(a, b...)
}

func (R *RClass) Parent() BaseObject {
	return R.superClass
}

func (R *RClass) SetParent(_ BaseObject) {}

// Helpers

func MakeClass(
	Name string,
	NewFunc CallFuncType,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  NewEnvironment(),
		superClass: nil,
		NewFunc:    NewFunc,
	}
}

func MakeClassFromSuper(
	Name string,
	super *RClass,
	NewFunc CallFuncType,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  NewEnvironment(),
		superClass: super,
		NewFunc:    NewFunc,
	}
}
