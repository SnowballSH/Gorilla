package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]BaseObject)
	f := make(map[string]BaseObject)
	return &Environment{Store: s, Free: f, outer: nil}
}

type Environment struct {
	Store map[string]BaseObject
	Free  map[string]BaseObject
	outer *Environment
}

func (e *Environment) Get(name string) (BaseObject, bool) {
	obj, ok := e.Free[name]
	if !ok {
		obj, ok = e.Store[name]
	}
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	if !ok {
		obj, ok = GlobalBuiltins[name]
	}
	return obj, ok
}
func (e *Environment) Set(name string, val BaseObject, macro bool) BaseObject {
	if macro {
		if _, ok := e.Free[name]; ok {
			e.Free[name] = val
		} else if _, ok = e.Store[name]; ok {
			e.Store[name] = val
		} else {
			e.Free[name] = val
		}
	} else {
		e.Store[name] = val
	}
	return val
}
