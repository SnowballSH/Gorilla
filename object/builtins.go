package object

var (
	BaseObjectBuiltins map[string]BaseObject
	IntegerBuiltins    map[string]BaseObject
)

func init() {
	BaseObjectBuiltins = map[string]BaseObject{
		"eq": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				return NewBool(self.Inspect() == args[0].(*Object).Inspect(), line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				return NewBool(self.Inspect() != args[0].(*Object).Inspect(), line)
			},
			[][]string{
				{ANY},
			},
		),
	}
	IntegerBuiltins = map[string]BaseObject{
		"add": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)+args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"sub": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)-args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"mul": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				return NewInteger(self.Value().(int)*args[0].(*Object).Value().(int), line)
			},
			[][]string{
				{INTEGER},
			},
		),
		"div": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
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
			func(self *Object, args []BaseObject, line int) BaseObject {
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
			func(self *Object, args []BaseObject, line int) BaseObject {
				v, ok := args[0].(*Object).Value().(int)
				if !ok {
					return NewBool(self.Inspect() != args[0].Inspect(), line)
				}
				return NewBool(self.Value().(int) == v, line)
			},
			[][]string{
				{ANY},
			},
		),
		"neq": NewBuiltinFunction(
			func(self *Object, args []BaseObject, line int) BaseObject {
				v, ok := args[0].(*Object).Value().(int)
				if !ok {
					return NewBool(self.Inspect() != args[0].Inspect(), line)
				}
				return NewBool(self.Value().(int) != v, line)
			},
			[][]string{
				{ANY},
			},
		),
	}
}
