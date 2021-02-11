package object

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

	IntMethods()
	FloatMethods()

	BooleanBuiltins = map[string]BaseObject{
		"not": NewBuiltinFunction(
			func(self *Object, env *Environment, args []BaseObject, line int) BaseObject {
				return NewBool(!self.Value().(bool), line)
			},
			[][]string{},
		),
	}

	StingMethods()

	ArrayMethods()
	HashMethod()

	GlobalFunctions()
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
