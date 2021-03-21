package runtime

import "sort"

type storeType map[string]BaseObject

func NewEnvironment() *environment {
	s := make(storeType)
	return &environment{store: s}
}

func NewEnvironmentWithStore(store storeType) *environment {
	return &environment{store: store}
}

type environment struct {
	store storeType
}

func (e *environment) get(name string) (BaseObject, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *environment) set(name string, val BaseObject) BaseObject {
	e.store[name] = val
	return val
}

func (e *environment) names() []string {
	var keys []string
	for key := range e.store {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (e *environment) copy() *environment {
	newEnv := make(storeType)
	for key, value := range e.store {
		newEnv[key] = value
	}
	return &environment{store: newEnv}
}
