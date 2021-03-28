package runtime

import "sort"

type storeType map[string]BaseObject

func NewEnvironment() *Environment {
	s := make(storeType)
	return &Environment{Store: s}
}

func NewEnvironmentWithStore(store storeType) *Environment {
	return &Environment{Store: store}
}

type Environment struct {
	Store storeType
}

func (e *Environment) Get(name string) (BaseObject, bool) {
	obj, ok := e.Store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val BaseObject) BaseObject {
	e.Store[name] = val
	return val
}

func (e *Environment) Names() []string {
	var keys []string
	for key := range e.Store {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (e *Environment) Copy() *Environment {
	newEnv := make(storeType)
	for key, value := range e.Store {
		newEnv[key] = value
	}
	return &Environment{Store: newEnv}
}
