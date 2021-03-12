package runtime

import "sort"

func newEnvironment() *environment {
	s := make(map[string]BaseObject)
	return &environment{store: s}
}

type environment struct {
	store map[string]BaseObject
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
	newEnv := make(map[string]BaseObject)
	for key, value := range e.store {
		newEnv[key] = value
	}
	return &environment{store: newEnv}
}
