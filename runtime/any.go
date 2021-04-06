package runtime

// All Gorilla object should inherit from this class
var AnyClass *RClass

var anyIns *Environment

func makeAnyIns() {
	anyIns = NewEnvironmentWithStore(map[string]BaseObject{
		"==": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			if self.Parent().EqualTo(k) {
				return NewInteger(1), nil
			} else {
				return NewInteger(0), nil
			}
		}),
		"!=": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			if !self.Parent().EqualTo(k) {
				return NewInteger(1), nil
			} else {
				return NewInteger(0), nil
			}
		}),
	})

	AnyClass = MakeClass("Any", NotCallable, anyIns)
}
