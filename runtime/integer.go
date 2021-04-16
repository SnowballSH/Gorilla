package runtime

import "fmt"

var IntegerClass *RClass

var intIns *Environment

func makeIntIns() {
	intIns = NewEnvironmentWithStore(map[string]BaseObject{
		// Integer + other: Integer -> Integer
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := k
			//return NewInteger(left + right), nil
			return NewInteger(add(left, right)), nil
		}),

		// Integer - other: Integer -> Integer
		"-": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}
			left := self.Parent().(*Object).InternalValue.(int64)
			right := k
			//return NewInteger(left - right), nil
			return NewInteger(sub(left, right)), nil
		}),

		// Integer * other: Integer -> Integer
		"*": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := k
			//return NewInteger(left * right), nil
			return NewInteger(mul(left, right)), nil
		}),

		// Integer / other: Integer -> Integer
		"/": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := k

			if right == 0 {
				return nil, fmt.Errorf("%d / %d: Division by 0", left, right)
			}

			return NewInteger(left / right), nil
		}),

		// Integer % other: Integer -> Integer
		"%": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := k

			if right == 0 {
				return nil, fmt.Errorf("%d %s %d: Division by 0", left, "%", right)
			}

			return NewInteger(left % right), nil
		}),

		// -Integer -> Integer
		"-@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return NewInteger(-left), nil
		}),

		// +Integer -> Integer
		"+@": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return NewInteger(+left), nil
		}),

		// Integer < other: Integer -> Bool
		"<": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			return fromBool(left < k), nil
		}),

		// Integer > other: Integer -> Bool
		">": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			return fromBool(left > k), nil
		}),

		// Integer <= other: Integer -> Bool
		"<=": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			return fromBool(left <= k), nil
		}),

		// Integer >= other: Integer -> Bool
		">=": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			k, err := getInteger(args, 0)
			if err != nil {
				return nil, err
			}

			return fromBool(left >= k), nil
		}),

		// nonz returns whether self is 0
		// Integer.nonz() -> Bool
		"nonz": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			return fromBool(left != 0), nil
		}),

		// times repeats a lambda self times
		// Integer.times(fn: Lambda) -> Null
		"times": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			fn, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			for i := int64(0); i < left; i++ {
				_, err = fn.Call(fn)
				if err != nil {
					return nil, err
				}
			}
			return Null, nil
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

func getInteger(args []BaseObject, index int) (int64, error) {
	k, err := getElement(args, index)
	if err != nil {
		return 0, err
	}
	ro, err := GorillaToInteger(k)
	if err != nil {
		return 0, err
	}
	return ro.InternalValue.(int64), nil
}
