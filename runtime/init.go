package runtime

import "fmt"

func init() {
	intIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			ro, err := GorillaToInteger(args[0])
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left + right), nil
		}),
		"-": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			ro, err := GorillaToInteger(args[0])
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left - right), nil
		}),
		"*": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {

			ro, err := GorillaToInteger(args[0])
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(int64)
			right := ro.InternalValue.(int64)
			return NewInteger(left * right), nil
		}),
		"/": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			ro, err := GorillaToInteger(args[0])
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
	})

	stringIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {

			ro, err := GorillaToString(args[0])
			if err != nil {
				return nil, err
			}

			left := self.Parent().(*Object).InternalValue.(string)
			right := ro.InternalValue.(string)
			return NewString(left + right), nil
		}),
	})

	GorillaToInteger = func(x BaseObject) (*Object, error) {
		o, ok := x.(*Object)
		if !ok {
			return nil, fmt.Errorf("cannot convert non-object to Integer")
		}

		switch o.Class() {
		case IntegerClass:
			return o, nil
		default:
			return nil, fmt.Errorf("cannot convert %s to Integer", o.Class().Name)
		}
	}

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
