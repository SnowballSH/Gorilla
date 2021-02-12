package stdlib

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/object"
	"io"
)

// File System
var FSNamespace = object.NewNameSpace("fs", map[string]object.BaseObject{})

const (
	WriterType = "file.Writer"
)

func fwf(self object.BaseObject) string {
	return fmt.Sprintf("file.Writer '%p'", self.Value())
}

func NewWriter(value io.Writer, line int) object.BaseObject {
	return object.NewObject(
		WriterType,
		value,
		fwf,
		fwf,
		line,
		map[string]object.BaseObject{
			"write": object.NewBuiltinFunction(
				func(self *object.Object, env *object.Environment, args []object.BaseObject, line int) object.BaseObject {
					_, e := self.Value().(io.Writer).Write([]byte(args[0].Value().(string)))
					if e != nil {
						return object.NewError(e.Error(), line)
					}
					return object.NULLOBJ
				},
				[][]string{{object.STRING}}),
		},
		nil,
		nil,
	)
}
