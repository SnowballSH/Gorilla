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

func NewLambda(params []string, bytecode []byte, oldVm *VM) *Object {
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
			if len(params) != len(args) {
				return nil, fmt.Errorf("expected %d arguments, got %d", len(params), len(args))
			}
			env := oldVm.Environment.Copy()
			for i, name := range params {
				env.Set(name, args[i])
			}
			vm := NewVMWithStore(bytecode, env)
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
