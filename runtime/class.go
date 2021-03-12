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

func (R *RClass) instanceVariables() *environment {
	return R.Instances
}

func (R *RClass) setInstanceVariables(e *environment) {
	R.Instances = e
}

func (R *RClass) isTruthy() bool {
	return true
}

func (R *RClass) equalTo(object BaseObject) bool {
	return object == R
}

// Helpers

func MakeClass(
	Name string,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  newEnvironment(),
		superClass: nil,
	}
}

func MakeClassFromSuper(
	Name string,
	super *RClass,
) *RClass {
	return &RClass{
		Name:       Name,
		Instances:  newEnvironment(),
		superClass: super,
	}
}
