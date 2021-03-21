package runtime

// Defines a Class, or Type
type RClass struct {
	Name       string
	Instances  *environment
	superClass *RClass
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

func (R *RClass) InstanceVariables() *environment {
	return R.Instances
}

func (R *RClass) SetInstanceVariables(e *environment) {
	R.Instances = e
}

func (R *RClass) IsTruthy() bool {
	return true
}

func (R *RClass) EqualTo(object BaseObject) bool {
	return object == R
}

func (R *RClass) Call(a BaseObject, b ...BaseObject) (BaseObject, error) {
	return NotCallable(a, b...)
}

func (R *RClass) Parent() BaseObject {
	return R.superClass
}

func (R *RClass) SetParent(_ BaseObject) {}

// Helpers

func MakeClass(
	Name string,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  NewEnvironment(),
		superClass: nil,
	}
}

func MakeClassFromSuper(
	Name string,
	super *RClass,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  NewEnvironment(),
		superClass: super,
	}
}
