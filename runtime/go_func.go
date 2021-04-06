package runtime

import "fmt"

var GoFuncClass *RClass

func makeGoFuncIns() {
	GoFuncClass = MakeClassFromSuper("Native Function", AnyClass, NotCallable, goFuncIns)
}

var goFuncIns = NewEnvironment()

func NewGoFunc(function CallFuncType) *Object {
	return &Object{
		RClass:        GoFuncClass,
		InternalValue: nil,
		ToStringFunc: func(self *Object) string {
			return "Native Function"
		},
		InspectFunc: func(self *Object) string {
			return fmt.Sprintf("Native Function %p", self)
		},
		IsTruthyFunc: func(self *Object) bool {
			return true
		},
		EqualToFunc: func(self *Object, other BaseObject) bool {
			return self == other
		},
		CallFunc: func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			return function(self, args...)
		},
		ParentObj: nil,
	}
}
