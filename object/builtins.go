package object

import (
	"Gorilla/config"
	"fmt"
	"strings"
	"unicode/utf8"
)

func init() {
	// Method that every Object has
	BaseObjectBuiltins = map[string]BaseObject{
		"isTruthy": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(
					self.Value() != false &&
						self.Value() != nil &&
						self.Value() != 0 &&
						self.Value() != "",
					line)
			},
			[][]string{},
		),

		"eq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(self.Inspect() == args[0].(*Object).Inspect() && self.Type() == args[0].Type(), line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(self.Inspect() != args[0].(*Object).Inspect() || self.Type() != args[0].Type(), line)
			},
			[][]string{
				{ANY},
			},
		),

		"and": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res, res1, err := getTwoBool(self, args[0].(*Object), env, line)
				if err != nil {
					return err
				}
				return NewBool(res && res1, line)
			},
			[][]string{
				{ANY},
			},
		),

		"or": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res, res1, err := getTwoBool(self, args[0].(*Object), env, line)
				if err != nil {
					return err
				}
				return NewBool(res || res1, line)
			},
			[][]string{
				{ANY},
			},
		),

		"toString": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(self.Inspect(), line)
			},
			[][]string{},
		),
		"toDebugString": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(self.Debug(), line)
			},
			[][]string{},
		),
	}

	NULLOBJ = NewNull(0)

	IntegerBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)+args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"sub": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)-args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"mul": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)*args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"div": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v := args[0].(*Object).Value().(int)
				if v == 0 {
					return NewError("Integer division by Zero", line)
				}
				return NewInteger(self.Value().(int)/v, line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"mod": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v := args[0].(*Object).Value().(int)
				if v == 0 {
					return NewError("Integer division by Zero", line)
				}
				return NewInteger(self.Value().(int)%v, line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"eq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				v, ok := args[0].(*Object).Value().(int)
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
				v, ok := args[0].(*Object).Value().(int)
				if !ok {
					return NewBool(true, line)
				}
				return NewBool(self.Value().(int) != v, line)
			},
			[][]string{
				{ANY},
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
	}

	BooleanBuiltins = map[string]BaseObject{
		"not": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(!self.Value().(bool), line)
			},
			[][]string{},
		),
	}

	StringBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(self.Value().(string)+args[0].Inspect(), line)
			},
			[][]string{{ANY}},
		),

		"mul": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(strings.Repeat(self.Value().(string), args[0].Value().(int)), line)
			},
			[][]string{{INTEGER}},
		),

		"length": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(utf8.RuneCountInString(self.Value().(string)), line)
			},
			[][]string{},
		),
	}

	GlobalBuiltins = map[string]BaseObject{
		"print": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var k []string
				for _, s := range args {
					k = append(k, s.Inspect())
				}
				_, _ = fmt.Fprint(config.OUT, strings.Join(k, " "))
				return NULLOBJ
			},
			nil,
		),
		"println": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var k []string
				for _, s := range args {
					k = append(k, s.Inspect())
				}
				_, _ = fmt.Fprintln(config.OUT, strings.Join(k, " "))
				return NULLOBJ
			},
			nil,
		),
		"debug": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var k []string
				for _, s := range args {
					k = append(k, s.Debug())
				}
				_, _ = fmt.Fprintln(config.OUT, strings.Join(k, " "))
				if len(args) > 0 {
					return args[len(args)-1]
				}
				return NULLOBJ
			},
			nil,
		),
		"type": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(args[0].Type(), line)
			},
			[][]string{
				{ANY},
			},
		),
		"null": NULLOBJ,
	}
}

var (
	GlobalBuiltins     map[string]BaseObject
	BaseObjectBuiltins map[string]BaseObject
	IntegerBuiltins    map[string]BaseObject
	BooleanBuiltins    map[string]BaseObject
	StringBuiltins     map[string]BaseObject

	NULLOBJ BaseObject
)

func getTwoBool(self *Object, other *Object, env *Environment, line int) (bool, bool, BaseObject) {
	fn, e := self.FindMethod("isTruthy")
	if e != nil {
		return false, false, e
	}

	res := fn.Call(env, self, []BaseObject{}, line)
	if res.Type() == ERROR {
		return false, false, res
	}
	if res.Type() != BOOLEAN {
		return false, false, NewError("isTruthy() Method expected to return Boolean", line)
	}

	fn, e = other.FindMethod("isTruthy")
	if e != nil {
		return false, false, e
	}

	res2 := fn.Call(env, other, []BaseObject{}, line)
	if res2.Type() == ERROR {
		return false, false, res2
	}
	if res2.Type() != BOOLEAN {
		return false, false, NewError("isTruthy() Method expected to return Boolean", line)
	}

	return res.Value().(bool), res2.Value().(bool), nil
}

func GetOneTruthy(self *Object, env *Environment, line int) (bool, BaseObject) {
	fn, e := self.FindMethod("isTruthy")
	if e != nil {
		return false, e
	}

	res := fn.Call(env, self, []BaseObject{}, line)
	if res.Type() == ERROR {
		return false, res
	}
	if res.Type() != BOOLEAN {
		return false, NewError("isTruthy() Method expected to return Boolean", line)
	}

	return res.Value().(bool), nil
}
