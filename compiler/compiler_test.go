package compiler

import (
	"Gorilla/code"
	"Gorilla/lexer"
	"Gorilla/object"
	"Gorilla/parser"
	"strconv"
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
		var m, n []string
		for _, v := range bytecodes {
			m = append(m, strconv.Itoa(int(v)))
		}
		for _, v := range _code {
			n = append(n, strconv.Itoa(int(v)))
		}
		t.Fatalf("Bytecode Length not same, Expected %d, got %d\nExpected Bytecode: %s, Got: %s",
			len(bytecodes), len(_code), m, n)
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

func Test2(t *testing.T) {
	assert(
		t,
		"i = 5; while i i = i - 1",
		[]code.Opcode{
			code.LoadConstant, // Load 5                        0
			code.SetVar,       // Set i to 5                    1
			code.Pop,          // Pop                           2
			code.GetVar,       // Get Variable 'i'              3
			code.JumpFalse,    // Jump                          4
			code.GetVar,       // Get Variable 'i'              5
			code.Method,       // Find i's "sub"                6
			code.LoadConstant, // Load 1                        7
			code.Call,         // Call i's "sub"                8
			code.SetVar,       // Set i to i - 1                9
			code.Pop,          // Pop                           10
			code.Jump,         // Jump Back                     11
			code.LoadConstant, // Load NULL                     12
			code.Pop,          // Pop                           13
		},
		[]object.Message{
			object.NewMessage(0),     // Load 5            0
			object.NewMessage("i"),   //                   1
			object.NewMessage("i"),   //                   2
			object.NewMessage(0),     //                   3
			object.NewMessage(11),    // 11 -> 12          4
			object.NewMessage(15),    // 15                5
			object.NewMessage("i"),   //                   6
			object.NewMessage(0),     //                   7
			object.NewMessage("sub"), //                   8
			object.NewMessage(1),     //                   9
			object.NewMessage(0),     //                   10
			object.NewMessage(1),     //                   11
			object.NewMessage("i"),   //                   12
			object.NewMessage(2),     // 2 -> 3            13
			object.NewMessage(2),     // 2                 14
			object.NewMessage(2),     // NULL              15
		},
		[]object.BaseObject{
			object.NewInteger(5, 0),
			object.NewInteger(1, 0),
			object.NULLOBJ,
		},
	)
}
