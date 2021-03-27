package runtime

import "sort"

type storeType map[string]BaseObject

func NewEnvironment() *Environment {
	s := make(storeType)
	return &Environment{store: s}
}

func NewEnvironmentWithStore(store storeType) *Environment {
	return &Environment{store: store}
}

type Environment struct {
	store storeType
}

func (e *Environment) get(name string) (BaseObject, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) set(name string, val BaseObject) BaseObject {
	e.store[name] = val
	return val
}

func (e *Environment) names() []string {
	var keys []string
	for key := range e.store {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (e *Environment) copy() *Environment {
	newEnv := make(storeType)
	for key, value := range e.store {
		newEnv[key] = value
	}
	return &Environment{store: newEnv}
}
