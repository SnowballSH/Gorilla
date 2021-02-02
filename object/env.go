package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]BaseObject)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]BaseObject
	outer *Environment
}

func (e *Environment) Get(name string) (BaseObject, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}
func (e *Environment) Set(name string, val BaseObject) BaseObject {
	e.store[name] = val
	return val
}
