package vm

import (
	"Gorilla/code"
	"Gorilla/object"
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
}

func NewFrame(bytecodes []code.Opcode, constants []object.BaseObject, messages []object.Message) *Frame {
	return &Frame{
		Instructions: bytecodes,
		Constants:    constants,
		Messages:     messages,
		Stack:        []object.BaseObject{},
		ip:           0,
		mp:           0,
		LastPopped:   nil,
		Env:          object.NewEnvironment(),
		LastFrame:    nil,
	}
}