package object

import (
	"fmt"
	"math"
)

func IntMethods() {
	IntegerBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				result := float64(self.Value().(int)) + otherv

				if float64(int(result)) == result {
					return NewInteger(int(result), line)
				}

				return NewFloat(result, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"sub": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				result := float64(self.Value().(int)) - otherv

				if float64(int(result)) == result {
					return NewInteger(int(result), line)
				}

				return NewFloat(result, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"mul": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				result := float64(self.Value().(int)) * otherv

				if float64(int(result)) == result {
					return NewInteger(int(result), line)
				}

				return NewFloat(result, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"div": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				if otherv == 0 {
					return NewError("Integer division by Zero", line)
				}

				result := float64(self.Value().(int)) / otherv

				if float64(int(result)) == result {
					return NewInteger(int(result), line)
				}

				return NewFloat(result, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"mod": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				if otherv == 0 {
					return NewError("Integer modulo by Zero", line)
				}

				result := math.Mod(float64(self.Value().(int)), otherv)

				if float64(int(result)) == result {
					return NewInteger(int(result), line)
				}

				return NewFloat(result, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),

		"to": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewIntRange(self.Value().(int), args[0].Value().(int))
			},
			[][]string{
				{INTEGER},
			},
		),

		"eq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v, ok := args[0].Value().(int)
				if !ok {
					return NewBool(false, line)
				}
				return NewBool(self.Value().(int) == v, line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v, ok := args[0].Value().(int)
				if !ok {
					return NewBool(true, line)
				}
				return NewBool(self.Value().(int) != v, line)
			},
			[][]string{
				{ANY},
			},
		),
		"lt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(float64(self.Value().(int)) < otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"gt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(float64(self.Value().(int)) > otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"lteq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(float64(self.Value().(int)) <= otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"gteq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(float64(self.Value().(int)) >= otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"pos": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(+self.Value().(int), line)
			},
			[][]string{},
		),
		"neg": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(-self.Value().(int), line)
			},
			[][]string{},
		),

		"chr": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(string(rune(self.Value().(int))), line)
			},
			[][]string{},
		),
		"toFloat": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(float64(self.Value().(int)), line)
			},
			[][]string{},
		),

		"pow": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewFloat(math.Pow(float64(self.Value().(int)), otherv), line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),

		"abs": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(int(math.Abs(float64(self.Value().(int)))), line)
			},
			[][]string{},
		),
		"log": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(math.Log(float64(self.Value().(int))), line)
			},
			[][]string{},
		),

		"times": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				fnv := args[0].Value().(*FunctionValue)
				fn := args[0]
				amountParams := len(fnv.Params)
				if amountParams == 0 {
					for i := 0; i < self.Value().(int); i++ {
						res := fn.Call(env, nil, []BaseObject{}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else if amountParams == 1 {
					for i := 0; i < self.Value().(int); i++ {
						res := fn.Call(env, nil, []BaseObject{NewInteger(i, line)}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else {
					return NewError(fmt.Sprintf("Integer.times function expects a macro with 0 or 1 parameters, got %d", amountParams), line)
				}
				return self
			},
			[][]string{{MACRO}},
		),
	}
}

func FloatMethods() {
	FloatBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				return NewFloat(self.Value().(float64)+otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"sub": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				return NewFloat(self.Value().(float64)-otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"mul": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}
				return NewFloat(self.Value().(float64)*otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"div": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				if otherv == 0 {
					return NewError("Float division by Zero", line)
				}

				return NewFloat(self.Value().(float64)/otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"mod": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				if otherv == 0 {
					return NewError("Float division by Zero", line)
				}

				return NewFloat(math.Mod(self.Value().(float64), otherv), line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"lt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(self.Value().(float64) < otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"gt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(self.Value().(float64) > otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"lteq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(self.Value().(float64) <= otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"gteq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewBool(self.Value().(float64) >= otherv, line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),
		"eq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v, ok := args[0].Value().(float64)
				if !ok {
					return NewBool(false, line)
				}
				return NewBool(self.Value().(float64) == v, line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v, ok := args[0].Value().(float64)
				if !ok {
					return NewBool(false, line)
				}
				return NewBool(self.Value().(float64) == v, line)
			},
			[][]string{
				{ANY},
			},
		),
		"pos": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(+self.Value().(float64), line)
			},
			[][]string{},
		),
		"neg": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(-self.Value().(float64), line)
			},
			[][]string{},
		),

		"toInt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res := int(self.Value().(float64))
				return NewInteger(res, line)
			},
			[][]string{},
		),

		"round": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res := math.Round(self.Value().(float64))
				return NewInteger(int(res), line)
			},
			[][]string{},
		),

		"floor": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res := math.Floor(self.Value().(float64))
				return NewInteger(int(res), line)
			},
			[][]string{},
		),

		"ceil": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res := math.Ceil(self.Value().(float64))
				return NewInteger(int(res), line)
			},
			[][]string{},
		),

		"pow": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var otherv float64
				if args[0].Type() == FLOAT {
					otherv = args[0].Value().(float64)
				} else if args[0].Type() == INTEGER {
					otherv = float64(args[0].Value().(int))
				}

				return NewFloat(math.Pow(self.Value().(float64), otherv), line)
			},
			[][]string{
				{FLOAT, INTEGER},
			},
		),

		"abs": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(math.Abs(self.Value().(float64)), line)
			},
			[][]string{},
		),
		"log": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewFloat(math.Log(self.Value().(float64)), line)
			},
			[][]string{},
		),
	}
}
