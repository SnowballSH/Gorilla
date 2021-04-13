package runtime

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/config"
	"strings"
)

//
var ReservedKW []string

// Global is the global runtime storage
var Global *Environment

func makeGlobal() {
	Global = NewEnvironmentWithStore(map[string]BaseObject{
		"$VERSION": NewString(config.VERSION),

		"null":  Null,
		"true":  GorillaTrue,
		"false": GorillaFalse,

		"print": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			var w []string
			for _, x := range args {
				w = append(w, x.ToString())
			}
			jw := strings.Join(w, " ")
			fmt.Println(jw)
			return NewString(jw), nil
		}),

		"str": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			k, err := getElement(args, 0)
			if err != nil {
				return nil, err
			}
			return NewString(k.ToString()), nil
		}),

		"Integer": IntegerClass,
	})

	ReservedKW = Global.Names()
}
