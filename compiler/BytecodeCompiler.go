package compiler

import (
	"Gorilla/ast"
	"Gorilla/code"
	"Gorilla/object"
)

type BytecodeCompiler struct {
	Constants []object.BaseObject
	Bytecodes []code.Opcode
	Messages  []object.Message
}

func NewBytecodeCompiler() *BytecodeCompiler {
	return &BytecodeCompiler{}
}

func (c *BytecodeCompiler) addConstant(obj object.BaseObject) int {
	c.Constants = append(c.Constants, obj)
	return len(c.Constants) - 1
}

func (c *BytecodeCompiler) addMessage(val interface{}) {
	c.Messages = append(c.Messages, object.NewMessage(val))
}

func (c *BytecodeCompiler) emit(co code.Opcode) {
	c.Bytecodes = append(c.Bytecodes, co)
}

func (c *BytecodeCompiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.Pop)

	case *ast.IntegerLiteral:
		c.addMessage(c.addConstant(object.NewInteger(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	default:
		panic("Node not supported: " + node.String())
	}
	return nil
}
