package runtime

import "fmt"

var IntegerClass *RClass

var intIns *Environment

func makeIntIns() {
	intIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToInteger(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left + right), nil
		}),

		"-": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToInteger(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left - right), nil
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

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left * right), nil
		}),

		"/": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToInteger(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)

			if right == 0 {
				return nil, fmt.Errorf("%d / %d: Division by 0", left, right)
			}

			return NewInteger(left / right), nil
		}),

		"%": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			ro, err := GorillaToInteger(k)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)

			if right == 0 {
				return nil, fmt.Errorf("%d %s %d: Division by 0", left, "%", right)
			}

			return NewInteger(left % right), nil
		}),

		"-@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return NewInteger(-left), nil
		}),

		"+@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return NewInteger(+left), nil
		}),

		"nonz": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return fromBool(left != 0), nil
		}),
	})

	IntegerClass = MakeClassFromSuper("Integer", NumericClass,
		func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			return GorillaToInteger(args[0])
		}, intIns)
}

var GorillaToInteger ConvertFuncType

func makeGorillaToInteger() {
	GorillaToInteger = func(x BaseObject) (*Object, error) {
		o, ok := x.(*Object)
		if !ok {
			return nil, fmt.Errorf("cannot convert non-object to Integer")
		}

		switch o.Class() {
		case IntegerClass:
			return o, nil
		case BoolClass:
			if o.InternalValue == true {
				return NewInteger(1), nil
			} else {
				return NewInteger(0), nil
			}
		default:
			return nil, fmt.Errorf("cannot convert %s to Integer", o.Class().Name)
		}
	}
}

func NewInteger(value int64) *Object {
	return &Object{
		RClass:        IntegerClass,
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
