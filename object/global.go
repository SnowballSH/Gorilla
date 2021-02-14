package object

import (
	"bufio"
	"fmt"
	"github.com/SnowballSH/Gorilla/config"
	"io"
	"os"
	"strings"
)

func GlobalFunctions() {
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
		"input": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var k []string
				for _, s := range args {
					k = append(k, s.Inspect())
				}
				_, _ = fmt.Fprint(config.OUT, strings.Join(k, " "))

				buffer := bufio.NewReader(os.Stdin)

				lineC, _, err := buffer.ReadLine()
				if err != nil && err != io.EOF {
					return NewError("EOF When getting input", line+1)
				}
				return NewString(string(lineC), line)
			},
			nil,
		),
		"exit": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				os.Exit(args[0].Value().(int))
				return NULLOBJ
			},
			[][]string{{INTEGER}},
		),

		"hash": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				value, ok := HashObject(args[0])
				if !ok {
					return NewError(fmt.Sprintf("Type '%s' is not hashable", args[0].Type()), line)
				}
				return NewFloat(float64(value.HashedKey), line)
			},
			[][]string{{ANY}},
		),

		// Experimental
		"makeObject": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				t := args[0].Value().(string)
				v := NewObject(
					"$"+t,
					nil,
					func(self BaseObject) string {
						return "Custom Object Type '" + t + "'"
					},
					func(self BaseObject) string {
						return "Custom Object Type '" + t + "'"
					},
					line,
					map[string]BaseObject{},
					nil,
					nil,
				)

				return v
			},
			[][]string{{STRING}},
		),

		"null":            NULLOBJ,
		"GORILLA_VERSION": NewString(config.VERSION, 0),
	}
}
