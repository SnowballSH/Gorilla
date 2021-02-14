package object

import "fmt"

func IntRangeMethods() {
	IntRangeBuiltins = map[string]BaseObject{
		"each": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				vv := self.Value().(*IntRangeValue)
				fnv := args[0].Value().(*FunctionValue)
				fn := args[0]
				amountParams := len(fnv.Params)
				if amountParams == 0 {
					for i := vv.start; i <= vv.end; i++ {
						res := fn.Call(env, nil, []BaseObject{}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else if amountParams == 1 {
					for i := vv.start; i <= vv.end; i++ {
						res := fn.Call(env, nil, []BaseObject{NewInteger(i, line)}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else {
					return NewError(fmt.Sprintf("Array.each function expects a macro with 0 or 1 parameters, got %d", amountParams), line)
				}
				return self
			},
			[][]string{{MACRO}},
		),
		"toArray": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var res []BaseObject
				vv := self.Value().(*IntRangeValue)
				for i := vv.start; i <= vv.end; i++ {
					res = append(res, NewInteger(i, line))
				}
				return NewArray(res, line)
			},
			[][]string{},
		),
	}
}
