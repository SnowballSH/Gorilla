package vm

import (
	"github.com/SnowballSH/Gorilla/code"
	"github.com/SnowballSH/Gorilla/object"
)

type Frame struct {
	Constants []object.BaseObject
	Messages  []object.Message
	mp        int

	Instructions []code.Opcode
	ip           int

	Stack []object.BaseObject

	LastPopped object.BaseObject

	Env *object.Environment

	LastFrame *Frame
	Macro     bool
}

func NewFrame(bytecodes []code.Opcode, constants []object.BaseObject, messages []object.Message) *Frame {
	e := object.NewEnvironment()
	f := &Frame{
		Instructions: bytecodes,
		Constants:    constants,
		Messages:     messages,
		Stack:        []object.BaseObject{},
		ip:           0,
		mp:           0,
		LastPopped:   nil,
		Env:          e,
		LastFrame:    nil,
		Macro:        false,
	}

	return f
}
