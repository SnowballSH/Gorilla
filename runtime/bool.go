package runtime

var BoolClass *RClass

var GorillaTrue BaseObject
var GorillaFalse BaseObject

var boolIns *Environment

var fromBool func(b bool) BaseObject

func makeBoolIns() {
	boolIns = NewEnvironmentWithStore(map[string]BaseObject{
		"!@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			return fromBool(!self.Parent().IsTruthy()), nil
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

	fromBool = func(b bool) BaseObject {
		if b {
			return GorillaTrue
		} else {
			return GorillaFalse
		}
	}
}
