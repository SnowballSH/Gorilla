package object

import (
	"Gorilla/config"
	"fmt"
	"strings"
)

var (
	GlobalBuiltins     map[string]BaseObject
	BaseObjectBuiltins map[string]BaseObject
	IntegerBuiltins    map[string]BaseObject

	NULLOBJ BaseObject = NewNull(0)
)

func init() {
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
		"null": NULLOBJ,
	}

	BaseObjectBuiltins = map[string]BaseObject{
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
	}

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
	}
}
