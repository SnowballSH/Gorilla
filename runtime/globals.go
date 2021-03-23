package runtime

import (
	"fmt"
	"strings"
)

var Global = NewEnvironmentWithStore(map[string]BaseObject{
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
