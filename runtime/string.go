package runtime

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

var StringClass *RClass

var stringIns *Environment

func makeStringIns() {
	stringIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToString(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(string)
			right := ro.InternalValue.(string)
			return NewString(left + right), nil
		}),

		"*": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToInteger(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(string)
			right := ro.InternalValue.(int64)
			return NewString(strings.Repeat(left, int(right))), nil
		}),
	})

	StringClass = MakeClassFromSuper("String", NumericClass,
		func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			return GorillaToString(args[0])
		}, stringIns)
}

var GorillaToString ConvertFuncType

func makeGorillaToString() {
	GorillaToString = func(x BaseObject) (*Object, error) {
		o, ok := x.(*Object)
		if !ok {
			return nil, fmt.Errorf("cannot convert non-object to String")
		}

		switch o.Class() {
		case StringClass:
			return o, nil
		default:
			return nil, fmt.Errorf("cannot convert %s to String", o.Class().Name)
		}
	}
}

func NewString(value string) *Object {
	return &Object{
		RClass:        StringClass,
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
