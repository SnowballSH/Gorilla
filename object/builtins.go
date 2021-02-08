package object

import (
	"Gorilla/config"
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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
				return NewBool(self.Inspect() == args[0].Inspect() && self.Type() == args[0].Type(), line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(self.Inspect() != args[0].Inspect() || self.Type() != args[0].Type(), line)
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
	}

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

		"getIndex": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				k := self.InternalValue.(string)
				idx := args[0].Value().(int)
				if idx < 0 {
					idx = utf8.RuneCountInString(k) + idx
				}
				if utf8.RuneCountInString(k) <= idx || idx < 0 {
					return NewError(fmt.Sprintf("String Index %d out of range on length %d", args[0].Value().(int), len(k)), line)
				}
				return NewString(string([]rune(k)[idx]), line)
			},
			[][]string{
				{INTEGER},
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
	}

	ArrayBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				for _, v := range args[0].Value().([]BaseObject) {
					self.InternalValue = append(self.InternalValue.([]BaseObject), v)
				}
				return self
			},
			[][]string{
				{ARRAY},
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
				idx := args[0].Value().(int)
				if idx < 0 {
					idx = len(k) + idx
				}
				if len(k) <= idx || idx < 0 {
					return NewError(fmt.Sprintf("Array Index %d out of range on length %d", args[0].Value().(int), len(k)), line)
				}
				return k[idx]
			},
			[][]string{
				{INTEGER},
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
				for _, val := range self.Value().([]BaseObject) {
					if val.Type() == args[0].Type() && val.Inspect() == args[0].Inspect() {
						ok = true
						break
					}
				}
				return NewBool(ok, line)
			},
			[][]string{{ANY}},
		),
	}

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

		"null":            NULLOBJ,
		"GORILLA_VERSION": NewString(config.VERSION, 0),
	}
}

var (
	GlobalBuiltins     map[string]BaseObject
	BaseObjectBuiltins map[string]BaseObject
	IntegerBuiltins    map[string]BaseObject
	FloatBuiltins      map[string]BaseObject
	BooleanBuiltins    map[string]BaseObject
	StringBuiltins     map[string]BaseObject
	ArrayBuiltins      map[string]BaseObject
	HashBuiltins       map[string]BaseObject

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
