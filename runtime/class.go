package runtime

// Defines a Class, or Type
type RClass struct {
	Name      string
	Instances *Environment

	// For children
	InstanceVars *Environment
	superClass   *RClass

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
	return R.Instances.Get(s)
}

func (R *RClass) InstanceVariables() *Environment {
	return R.Instances
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

// GetInstance gets instance variables for its children
func (R *RClass) GetInstance(s string) (BaseObject, bool) {
	o, ok := R.InstanceVars.Get(s)
	if !ok && R.superClass != nil {
		o, ok = R.superClass.GetInstance(s)
	}
	return o, ok
}

// Helpers

func MakeClass(
	Name string,
	NewFunc CallFuncType,
	InstanceVars *Environment,
) *RClass {
	return &RClass{
		Name:         Name,
		Instances:    NewEnvironment(),
		InstanceVars: InstanceVars,
		superClass:   nil,
		NewFunc:      NewFunc,
	}
}

func MakeClassFromSuper(
	Name string,
	super *RClass,
	NewFunc CallFuncType,
	InstanceVars *Environment,
) *RClass {
	return &RClass{
		Name:         Name,
		Instances:    NewEnvironment(),
		InstanceVars: InstanceVars,
		superClass:   super,
		NewFunc:      NewFunc,
	}
}
