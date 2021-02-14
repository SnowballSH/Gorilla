package object

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/helper"
	"strconv"
	"strings"
	"unicode/utf8"
)

func StingMethods() {
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

		"getIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(string)
				w := utf8.RuneCountInString(k)

				if args[0].Type() == INTEGER {
					idx := args[0].Value().(int)
					if idx < 0 {
						idx += w
					}
					if w <= idx || idx < 0 {
						return NewError(fmt.Sprintf("String Index %d out of range on length %d", args[0].Value().(int), len(k)), line)
					}
					return NewString(string([]rune(k)[idx]), line)
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

					val := string([]rune(k)[start : end+1])

					if reverse {
						val = helper.ReverseString(val)
					}

					return NewString(val, line)
				}
			},
			[][]string{
				{INTEGER, INTRANGE},
			},
		),

		"length": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewInteger(utf8.RuneCountInString(self.Value().(string)), line)
			},
			[][]string{},
		),

		"toInt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res, err := strconv.Atoi(self.Value().(string))
				if err != nil {
					return NewError(fmt.Sprintf("'%s' is not a valid integer", self.Value().(string)), line)
				}
				return NewInteger(res, line)
			},
			[][]string{},
		),

		"toFloat": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				res, err := strconv.ParseFloat(self.Value().(string), 64)
				if err != nil {
					return NewError(fmt.Sprintf("'%s' is not a valid float", self.Value().(string)), line)
				}
				return NewFloat(res, line)
			},
			[][]string{},
		),

		"isInt": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				_, err := strconv.Atoi(self.Value().(string))
				return NewBool(err == nil, line)
			},
			[][]string{},
		),

		"isFloat": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				_, err := strconv.ParseFloat(self.Value().(string), 64)
				return NewBool(err == nil, line)
			},
			[][]string{},
		),

		"ords": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				r := []rune(self.Value().(string))
				rr := make([]BaseObject, len(r))
				for i, v := range r {
					rr[i] = NewInteger(int(v), line)
				}
				return NewArray(rr, line)
			},
			[][]string{},
		),

		"ord": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				r := []rune(self.Value().(string))
				if len(r) == 0 {
					return NewError(fmt.Sprintf("String is empty"), line)
				}
				return NewInteger(int(r[0]), line)
			},
			[][]string{},
		),

		"chars": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				var arr []BaseObject
				for _, v := range self.Value().(string) {
					arr = append(arr, NewString(string(v), line))
				}
				return NewArray(arr, line)
			},
			[][]string{},
		),

		"strip": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(strings.TrimSpace(self.Value().(string)), line)
			},
			[][]string{},
		),

		"cut": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(strings.Trim(self.Value().(string), args[0].Value().(string)), line)
			},
			[][]string{{STRING}},
		),

		"upper": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(strings.ToUpper(self.Value().(string)), line)
			},
			[][]string{},
		),

		"lower": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewString(strings.ToLower(self.Value().(string)), line)
			},
			[][]string{},
		),

		"has": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(strings.Contains(self.Value().(string), args[0].Value().(string)), line)
			},
			[][]string{{STRING}},
		),
	}
}
