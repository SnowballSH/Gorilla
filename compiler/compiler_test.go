package compiler

import (
	"Gorilla/code"
	"Gorilla/lexer"
	"Gorilla/object"
	"Gorilla/parser"
	"testing"
)

func assert(t *testing.T, code string, bytecodes []code.Opcode, messages []object.Message, constants []object.BaseObject) {
	cp := NewBytecodeCompiler()
	l := lexer.New(code)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("Parser Error:\n\t%s", p.Errors())
	}

	err := cp.Compile(program)
	if err != nil {
		t.Fatalf("Compiler Error:\n\t%s", err.Error())
	}

	_code := cp.Bytecodes
	_constants := cp.Constants
	_messages := cp.Messages

	if len(bytecodes) != len(_code) {
		t.Fatalf("Bytecode Length not same, Expected %d, got %d", len(bytecodes), len(_code))
	}
	for i, v := range bytecodes {
		if v != _code[i] {
			t.Fatalf("Bytecode #%d not same, Expected %d, got %d", i, v, _code[i])
		}
	}

	if len(constants) != len(_constants) {
		t.Fatalf("Constant Length not same, Expected %d, got %d", len(constants), len(_constants))
	}
	for i, v := range constants {
		if v.Inspect() != _constants[i].Inspect() {
			t.Fatalf("Constant #%d not same, Expected %s, got %s", i, v.Debug(), _constants[i].Debug())
		}
	}

	if len(messages) != len(_messages) {
		t.Fatalf("Message Length not same, Expected %d, got %d", len(messages), len(_messages))
	}
	for i, v := range messages {
		if v.MessageType() != _messages[i].MessageType() {
			t.Fatalf("Message #%d type not same, Expected %d, got %d", i, v.MessageType(), _messages[i].MessageType())
		}
		if v.MessageType() == object.INTMESSAGE {
			if v.(*object.IntMessage).Value != _messages[i].(*object.IntMessage).Value {
				t.Fatalf("Message #%d not same, Expected %d, got %d", i,
					v.(*object.IntMessage).Value, _messages[i].(*object.IntMessage).Value)
			}
		}
	}
}

func Test1(t *testing.T) {
	assert(
		t,
		"if false 1 else 2",
		[]code.Opcode{
			code.LoadConstant, // Load false                     0
			code.JumpFalse,    // Jump If false is false        1
			code.LoadConstant, // Load 1                        2
			code.Jump,         // Jump Out                      3
			code.LoadConstant, // Else, Load 2                  4
			code.Pop,          // Pop                           5
		},
		[]object.Message{
			object.NewMessage(0), // Load false             0
			object.NewMessage(3), // Jump to 3 -> 4        1
			object.NewMessage(6), // Jump Message to 6     2
			object.NewMessage(1), // Load 1                3
			object.NewMessage(4), // Jump to 4 -> 5        4
			object.NewMessage(7), // Jump Message to 7     5
			object.NewMessage(2), // Load 2                6
		},
		[]object.BaseObject{
			object.NewBool(false, 0),
			object.NewInteger(1, 0),
			object.NewInteger(2, 0),
		},
	)
}
