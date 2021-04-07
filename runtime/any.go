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
				return GorillaTrue, nil
			} else {
				return GorillaFalse, nil
			}
		}),
		"!=": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			if !self.Parent().EqualTo(k) {
				return GorillaTrue, nil
			} else {
				return GorillaFalse, nil
			}
		}),
	})

	AnyClass = MakeClass("Any", NotCallable, anyIns)
}
