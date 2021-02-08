package stdlib

import (
	"Gorilla/object"
)

var StandardLibrary = map[string]object.BaseObject{
	"math":   MathNamespace,
	"random": RandomNamespace,
	"http":   HttpNamespace,
	"fs":     FSNamespace,
}
