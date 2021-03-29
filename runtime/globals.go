package runtime

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/config"
	"strings"
)

// Global is the global runtime storage
var Global *Environment

func makeGlobal() {
	Global = NewEnvironmentWithStore(map[string]BaseObject{
		"$VERSION": NewString(config.VERSION),
		"null":     Null,

		"print": NewGoFunc(func(self BaseObject, args ...BaseObject) (BaseObject, error) {
			var w []string
			for _, x := range args {
				w = append(w, x.ToString())
			}
			jw := strings.Join(w, " ")
			fmt.Println(jw)
			return NewString(jw), nil
		}),
	})
}
