package stdlib

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/config"
	"github.com/SnowballSH/Gorilla/object"
	"hash/fnv"
	"net/http"
	"os"
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
				res := fnc.Call(env, nil, []object.BaseObject{
					NewWriter(writer, line),
					NewRequest(request, line),
				}, line)
				if res.Type() == object.ERROR {
					_, _ = fmt.Fprintln(config.OUT, " Runtime Error (while serving):\n\t"+res.Inspect())
					os.Exit(2)
				}
			}

			http.HandleFunc(args[0].Value().(string), fn)
			return object.NULLOBJ
		},
		[][]string{{object.STRING}, {object.FUNCTION}},
	),
})

const (
	RequestType = "http.Request"
)

func rqt(self object.BaseObject) string {
	return fmt.Sprintf("http.Reqeust '%p'", self.Value())
}

func NewRequest(value *http.Request, line int) object.BaseObject {
	header := map[object.HashKey]*object.HashValue{}
	for s, v := range value.Header {
		h := fnv.New64a()
		_, _ = h.Write([]byte(s))

		var bm []object.BaseObject
		for _, l := range v {
			bm = append(bm, object.NewString(l, line))
		}

		header[object.HashKey{Type: object.STRING, HashedKey: h.Sum64()}] = &object.HashValue{
			Key:   object.NewString(s, line),
			Value: object.NewArray(bm, line),
		}
	}

	headerobj := object.NewHash(header, line)

	return object.NewObject(
		RequestType,
		value,
		rqt,
		rqt,
		line,
		map[string]object.BaseObject{
			"method": object.NewString(value.Method, line),
			"proto":  object.NewString(value.Proto, line),
			"header": headerobj,
		},
		nil,
		nil,
	)
}
