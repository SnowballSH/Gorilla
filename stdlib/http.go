package stdlib

import (
	"Gorilla/object"
	"fmt"
	"net/http"
)

// Http server
var HttpNamespace = object.NewNameSpace("http", map[string]object.BaseObject{
	"serve": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			e := http.ListenAndServe(args[0].Value().(string), nil)
			if e != nil {
				return object.NewError(e.Error(), line)
			}
			return object.NULLOBJ
		},
		[][]string{{object.STRING}},
	),
	"get": object.NewBuiltinFunction(
		func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
			fnc := args[1]
			amount := len(fnc.Value().(*object.FunctionValue).Params)
			if amount != 2 {
				return object.NewError(fmt.Sprintf("'http.get' function requires a function with 2 parameters. got %d", amount), line)
			}
			fn := func(writer http.ResponseWriter, request *http.Request) {
				fnc.Call(env, nil, []object.BaseObject{
					NewWriter(writer, line),
					NewRequest(request, line),
				}, line)
			}

			http.HandleFunc(args[0].Value().(string), fn)
			return object.NULLOBJ
		},
		[][]string{{object.STRING}, {object.FUNCTION}},
	),
})

const (
	RequestType = "http.request"
)

func rqt(self object.BaseObject) string {
	return fmt.Sprintf("file.Writer '%p'", self.Value())
}

func NewRequest(value *http.Request, line int) object.BaseObject {
	return object.NewObject(
		RequestType,
		value,
		rqt,
		rqt,
		line,
		map[string]object.BaseObject{},
		nil,
		nil,
	)
}
