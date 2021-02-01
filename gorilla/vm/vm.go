package vm

import "../object"
import _ "github.com/alecthomas/participle"

type VM struct {
	constants []object.Object
}
