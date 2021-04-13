package runtime

import (
	"fmt"
)

var LambdaClass *RClass

var lambdaIns *Environment

func makeLambdaIns() {
	lambdaIns = NewEnvironment()

	LambdaClass = MakeClassFromSuper("Lambda Function", AnyClass, NotCallable, lambdaIns)
}

func NewLambda(bytecode []byte, oldVm *VM) *Object {
	return &Object{
		RClass:        LambdaClass,
		InternalValue: bytecode,
		ToStringFunc: func(self *Object) string {
			return "Lambda Function"
		},
		InspectFunc: func(self *Object) string {
			return fmt.Sprintf("Lambda Function %p", self)
		},
		IsTruthyFunc: func(self *Object) bool {
			return true
		},
		EqualToFunc: func(self *Object, other BaseObject) bool {
			return self == other
		},
		CallFunc: func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			vm := NewVMWithStore(bytecode, oldVm.Environment.Copy())
			vm.Run()
			k := vm.LastPopped
			if k == nil {
				k = Null
			}
			if vm.Error != nil {
				return k, fmt.Errorf(
					"Runtime Error in line %d:\n\n%s",
					vm.Error.Line+1,
					vm.Error.Message,
				)
			}
			return k, nil
		},
		ParentObj: nil,
	}
}
