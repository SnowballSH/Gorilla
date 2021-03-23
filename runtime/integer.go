package runtime

import "fmt"

var IntegerClass = MakeClassFromSuper("Integer", NumericClass,
	func(self BaseObject, args ...BaseObject) (BaseObject, error) {
		return GorillaToInteger(args[0])
	})

var intIns *environment

var GorillaToInteger ConvertFuncType

func NewInteger(value int64) *Object {
	return &Object{
		RClass:        IntegerClass,
		Instances:     intIns,
		InternalValue: value,
		ToStringFunc: func(self *Object) string {
			return fmt.Sprintf("%d", self.InternalValue.(int64))
		},
		InspectFunc: func(self *Object) string {
			return fmt.Sprintf("%d", self.InternalValue.(int64))
		},
		IsTruthyFunc: func(self *Object) bool {
			return self.InternalValue.(int64) != 0
		},
		EqualToFunc: func(self *Object, other BaseObject) bool {
			v, o := other.(*Object)
			if !o {
				return false
			}
			x, o := v.InternalValue.(int64)
			if !o {
				return false
			}
			return x == self.InternalValue.(int64)
		},
		CallFunc:  NotCallable,
		ParentObj: nil,
	}
}
