package runtime

var BoolClass *RClass

var GorillaTrue BaseObject
var GorillaFalse BaseObject

var boolIns *Environment

func makeBoolIns() {
	boolIns = NewEnvironmentWithStore(map[string]BaseObject{
		"!@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			if self.Parent().IsTruthy() {
				return GorillaFalse, nil
			} else {
				return GorillaTrue, nil
			}
		}),
	})

	BoolClass = MakeClassFromSuper("Bool", AnyClass, NotCallable, boolIns)

	GorillaTrue = &Object{
		RClass:        BoolClass,
		InternalValue: true,
		ToStringFunc: func(self *Object) string {
			return "true"
		},
		InspectFunc: func(self *Object) string {
			return "true"
		},
		IsTruthyFunc: func(self *Object) bool {
			return true
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

	GorillaFalse = &Object{
		RClass:        BoolClass,
		InternalValue: false,
		ToStringFunc: func(self *Object) string {
			return "false"
		},
		InspectFunc: func(self *Object) string {
			return "false"
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
}
