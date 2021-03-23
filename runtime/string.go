package runtime

import (
	"fmt"
	"unicode/utf8"
)

var StringClass = MakeClassFromSuper("String", NumericClass,
	func(self BaseObject, args ...BaseObject) (BaseObject, error) {
		return GorillaToString(args[0])
	})

var stringIns *environment

var GorillaToString ConvertFuncType

func NewString(value string) *Object {
	return &Object{
		RClass:        StringClass,
		Instances:     stringIns,
		InternalValue: value,
		ToStringFunc: func(self *Object) string {
			return fmt.Sprintf("%s", self.InternalValue.(string))
		},
		InspectFunc: func(self *Object) string {
			return fmt.Sprintf("'%s'", self.InternalValue.(string))
		},
		IsTruthyFunc: func(self *Object) bool {
			return utf8.RuneCountInString(self.InternalValue.(string)) != 0
		},
		EqualToFunc: func(self *Object, other BaseObject) bool {
			v, o := other.(*Object)
			if !o {
				return false
			}
			x, o := v.InternalValue.(string)
			if !o {
				return false
			}
			return x == self.InternalValue.(string)
		},
		CallFunc:  NotCallable,
		ParentObj: nil,
	}
}
