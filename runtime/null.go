package runtime

import (
	"fmt"
)

var NullClass = MakeClassFromSuper("Null", AnyClass,
	func(self BaseObject, args ...BaseObject) (BaseObject, error) {
		return nil, fmt.Errorf("cannot convert %s to null", args[0].Class().Name)
	}, nullIns)

var nullIns *Environment

func makeNullIns() {
	nullIns = NewEnvironmentWithStore(map[string]BaseObject{})
}

var Null = &Object{
	RClass:        NullClass,
	InternalValue: nil,
	ToStringFunc: func(self *Object) string {
		return "null"
	},
	InspectFunc: func(self *Object) string {
		return "null"
	},
	IsTruthyFunc: func(self *Object) bool {
		return false
	},
	EqualToFunc: func(self *Object, other BaseObject) bool {
		v, o := other.(*Object)
		if !o {
			return false
		}
		return self == v
	},
	CallFunc:  NotCallable,
	ParentObj: nil,
}
