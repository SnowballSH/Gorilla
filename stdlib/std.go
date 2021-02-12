package stdlib

import (
	"github.com/SnowballSH/Gorilla/object"
)

var StandardLibrary = map[string]object.BaseObject{
	"math":   MathNamespace,
	"random": RandomNamespace,
	"http":   HttpNamespace,
	"fs":     FSNamespace,
}
