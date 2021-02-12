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
		"createClass": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				hash := args[1].Value().(map[HashKey]*HashValue)
				v := NewHash(hash, line)
				v.TT = "$" + args[0].Value().(string)

				for _, vv := range hash {
					v.Methods[vv.Key.Inspect()] = vv.Value
				}

				return v
			},
			[][]string{{STRING}, {HASH}},
		),

		"null":            NULLOBJ,
		"GORILLA_VERSION": NewString(config.VERSION, 0),
	}
}
