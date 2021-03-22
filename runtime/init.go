package runtime

import "fmt"

func init() {
	intIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			right := args[0].(*Object).InternalValue.(int64)
			return NewInteger(left + right), nil
		}),
		"-": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			right := args[0].(*Object).InternalValue.(int64)
			return NewInteger(left - right), nil
		}),
		"*": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			right := args[0].(*Object).InternalValue.(int64)
			return NewInteger(left * right), nil
		}),
		"/": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(int64)
			right := args[0].(*Object).InternalValue.(int64)
			if right == 0 {
				return nil, fmt.Errorf("%d / %d: Division by 0", left, right)
			}
			return NewInteger(left / right), nil
		}),
	})

	stringIns = NewEnvironmentWithStore(map[string]BaseObject{
		"+": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			left := self.Parent().(*Object).InternalValue.(string)
			right := args[0].(*Object).InternalValue.(string)
			return NewString(left + right), nil
		}),
	})
}
