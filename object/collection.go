package object

import (
	"fmt"
	"strings"
)

func ArrayMethods() {
	ArrayBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				arr := CopyObject(self)
				for _, v := range args[0].Value().([]BaseObject) {
					arr.InternalValue = append(arr.InternalValue.([]BaseObject), v)
				}
				return arr
			},
			[][]string{
				{ARRAY},
			},
		),
		"mul": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				i := 0
				a := args[0].Value().(int)
				arr := NewArray([]BaseObject{}, line)
				for i < a {
					for _, v := range self.Value().([]BaseObject) {
						arr.InternalValue = append(arr.InternalValue.([]BaseObject), v)
					}
					i++
				}
				return arr
			},
			[][]string{
				{INTEGER},
			},
		),

		"push": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				self.InternalValue = append(self.InternalValue.([]BaseObject), args[0])
				return self
			},
			[][]string{
				{ANY},
			},
		),
		"pop": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.([]BaseObject)
				if len(k) < 1 {
					return NewError("Cannot pop empty list", line)
				}
				v := k[len(k)-1]
				self.InternalValue = k[:len(k)-1]
				return v
			},
			[][]string{},
		),
		"shift": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.([]BaseObject)
				if len(k) < 1 {
					return NewError("Cannot shift empty list", line)
				}
				v := k[0]
				self.InternalValue = k[1:]
				return v
			},
			[][]string{},
		),

		"getIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.([]BaseObject)
				w := len(k)

				if args[0].Type() == INTEGER {
					idx := args[0].Value().(int)
					if idx < 0 {
						idx += w
					}
					if w <= idx || idx < 0 {
						return NewError(fmt.Sprintf("Array Index %d out of range on length %d", args[0].Value().(int), len(k)), line)
					}
					return k[idx]
				} else /* INTRANGE */ {
					v := args[0].Value().(*IntRangeValue)
					start := v.start
					end := v.end

					if start < 0 {
						start += w
					}
					if end < 0 {
						end += w
					}

					if w <= start || start < 0 {
						return NewError(fmt.Sprintf("String Index %d out of range on length %d", v.start, len(k)), line)
					}
					if w <= end || end < 0 {
						return NewError(fmt.Sprintf("String Index %d out of range on length %d", v.end, len(k)), line)
					}

					var reverse = false

					if start > end {
						reverse = true
						start, end = end, start
					}

					val := k[start : end+1]

					if reverse {
						var a []BaseObject
						for _, v := range val {
							a = append(a, v)
						}
						for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
							a[i], a[j] = a[j], a[i]
						}
						val = a
					}

					return NewArray(val, line)
				}
			},
			[][]string{
				{INTEGER, INTRANGE},
			},
		),
		"setIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.([]BaseObject)
				idx := args[0].Value().(int)
				if idx < 0 {
					idx = len(k) + idx
				}
				if len(k) <= idx || idx < 0 {
					return NewError(fmt.Sprintf("Array Index %d out of range on length %d", args[0].Value().(int), len(k)), line)
				}
				k[idx] = args[1]
				self.InternalValue = k
				return self
			},
			[][]string{
				{INTEGER},
				{ANY},
			},
		),

		"length": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(len(self.Value().([]BaseObject)), line)
			},
			[][]string{},
		),
		"has": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				ok := false
				t := args[0].Type()
				is := args[0].Inspect()
				for _, val := range self.Value().([]BaseObject) {
					if val.Type() == t && val.Inspect() == is {
						ok = true
						break
					}
				}
				return NewBool(ok, line)
			},
			[][]string{{ANY}},
		),

		"join": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				js := ", "
				if len(args) > 1 {
					return NewError(
						fmt.Sprintf("Argument amount mismatch: Expected 0 or 1, got %d", len(args)),
						line,
					)
				}
				if len(args) == 1 {
					js = args[0].Inspect()
				}

				var ss []string
				for _, v := range self.Value().([]BaseObject) {
					ss = append(ss, v.Inspect())
				}

				return NewString(strings.Join(ss, js), line)
			},
			nil,
		),

		"map": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				fnv := args[0].Value().(*FunctionValue)
				fn := args[0]
				amountParams := len(fnv.Params)
				if amountParams == 0 {
					for i := range self.Value().([]BaseObject) {
						res := fn.Call(env, nil, []BaseObject{}, line)
						if res.Type() == ERROR {
							return res
						}
						self.Value().([]BaseObject)[i] = res
					}
				} else if amountParams == 1 {
					for i, val := range self.Value().([]BaseObject) {
						res := fn.Call(env, nil, []BaseObject{val}, line)
						if res.Type() == ERROR {
							return res
						}
						self.Value().([]BaseObject)[i] = res
					}
				} else {
					return NewError(fmt.Sprintf("Array.map function expects a function with 0 or 1 parameters, got %d", amountParams), line)
				}
				return self
			},
			[][]string{{FUNCTION}},
		),
		"each": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				fnv := args[0].Value().(*FunctionValue)
				fn := args[0]
				amountParams := len(fnv.Params)
				if amountParams == 0 {
					for range self.Value().([]BaseObject) {
						res := fn.Call(env, nil, []BaseObject{}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else if amountParams == 1 {
					for _, val := range self.Value().([]BaseObject) {
						res := fn.Call(env, nil, []BaseObject{val}, line)
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

		"reverse": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				val := self.Value().([]BaseObject)
				var a []BaseObject
				for _, v := range val {
					a = append(a, v)
				}
				for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
					a[i], a[j] = a[j], a[i]
				}
				return NewArray(a, line)
			},
			[][]string{},
		),
	}
}

func HashMethod() {
	HashBuiltins = map[string]BaseObject{
		"getIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(map[HashKey]*HashValue)
				key, ok := HashObject(args[0])
				if !ok {
					return NewError(fmt.Sprintf("Type '%s' is not hashable", args[0].Type()), line)
				}
				value, get := k[key]
				if !get {
					return NewError(fmt.Sprintf("Key not found: %s", args[0].Debug()), line)
				}

				return value.Value
			},
			[][]string{
				{ANY},
			},
		),
		"setIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(map[HashKey]*HashValue)
				key, ok := HashObject(args[0])
				if !ok {
					return NewError(fmt.Sprintf("Type '%s' is not hashable", args[0].Type()), line)
				}
				k[key] = &HashValue{
					Key:   args[0],
					Value: args[1],
				}
				self.InternalValue = k
				return self
			},
			[][]string{
				{ANY},
				{ANY},
			},
		),
		"values": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(map[HashKey]*HashValue)
				var values []BaseObject
				for _, v := range k {
					values = append(values, v.Value)
				}
				return NewArray(values, line)
			},
			[][]string{},
		),
		"keys": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(map[HashKey]*HashValue)
				var keys []BaseObject
				for _, v := range k {
					keys = append(keys, v.Key)
				}
				return NewArray(keys, line)
			},
			[][]string{},
		),
		"items": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(map[HashKey]*HashValue)
				var keys []BaseObject
				for _, v := range k {
					keys = append(keys, NewArray([]BaseObject{v.Key, v.Value}, line))
				}
				return NewArray(keys, line)
			},
			[][]string{},
		),
		"length": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(len(self.InternalValue.(map[HashKey]*HashValue)), line)
			},
			[][]string{},
		),

		"each": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				fnv := args[0].Value().(*FunctionValue)
				fn := args[0]
				amountParams := len(fnv.Params)
				if amountParams == 0 {
					for range self.Value().(map[HashKey]*HashValue) {
						res := fn.Call(env, nil, []BaseObject{}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else if amountParams == 1 {
					for _, val := range self.Value().(map[HashKey]*HashValue) {
						res := fn.Call(env, nil, []BaseObject{val.Key}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else if amountParams == 2 {
					for _, val := range self.Value().(map[HashKey]*HashValue) {
						res := fn.Call(env, nil, []BaseObject{val.Key, val.Value}, line)
						if res.Type() == ERROR {
							return res
						}
					}
				} else {
					return NewError(fmt.Sprintf("Hash.each function expects a macro with 0 to 2 parameters, got %d", amountParams), line)
				}
				return self
			},
			[][]string{{MACRO}},
		),
	}
}
