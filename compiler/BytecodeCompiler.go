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

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		c.addMessage(node.Name.Value)
		c.emit(code.SetVar)

	case *ast.AssignmentExpression:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		c.addMessage(node.Name.Value)
		c.emit(code.SetVar)

	case *ast.Identifier:
		c.addMessage(node.Value)
		c.addMessage(node.Token.Line)
		c.emit(code.GetVar)

	case *ast.IntegerLiteral:
		c.addMessage(c.addConstant(object.NewInteger(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.Boolean:
		c.addMessage(c.addConstant(object.NewBool(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.InfixExpression:
		name := ""
		switch node.Operator {
		case "+":
			name = "add"
		case "-":
			name = "sub"
		case "*":
			name = "mul"
		case "/":
			name = "div"
		case "%":
			name = "mod"
		case "==":
			name = "eq"
		case "!=":
			name = "neq"
		case ">":
			name = "gt"
		case "<":
			name = "lt"
		case ">=":
			name = "gteq"
		case "<=":
			name = "lteq"
		case "<-":
			name = "push"
		case "||":
			name = "or"
		case "&&":
			name = "and"
		default:
			panic("Operator not handled: " + node.Operator)
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		c.addMessage(name)
		c.emit(code.Method)

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		c.addMessage(node.Token.Line)
		c.addMessage(1)
		c.emit(code.Call)

	case *ast.GetAttr:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}
		c.addMessage(node.Name.String())
		c.emit(code.Method)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, k := range node.Arguments {
			err = c.Compile(k)
			if err != nil {
				return err
			}
		}

		c.addMessage(node.Token.Line)
		c.addMessage(len(node.Arguments))
		c.emit(code.Call)

	default:
		panic("Node not supported: " + node.String())
	}
	return nil
}
