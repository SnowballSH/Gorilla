package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]BaseObject)
	return &Environment{Store: s, outer: nil}
}

type Environment struct {
	Store map[string]BaseObject
	outer *Environment
}

func (e *Environment) Get(name string) (BaseObject, bool) {
	obj, ok := e.Store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	if !ok {
		obj, ok = GlobalBuiltins[name]
	}
	return obj, ok
}
func (e *Environment) Set(name string, val BaseObject) BaseObject {
	e.Store[name] = val
	return val
}
