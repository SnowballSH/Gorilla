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

func (c *BytecodeCompiler) addMessage(val interface{}) int {
	c.Messages = append(c.Messages, object.NewMessage(val))
	return len(c.Messages) - 1
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

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

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

	case *ast.StringLiteral:
		c.addMessage(c.addConstant(object.NewString(node.Value, node.Token.Line)))
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

	case *ast.PrefixExpression:
		name := ""
		switch node.Operator {
		case "!":
			name = "not"
		case "-":
			name = "neg"
		case "+":
			name = "pos"
		default:
			panic("Prefix Operator not handled: " + node.Operator)
		}

		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		c.addMessage(name)
		c.emit(code.Method)

		c.addMessage(node.Token.Line)
		c.addMessage(0)
		c.emit(code.Call)

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		c.emit(code.JumpFalse)
		m1 := c.addMessage(0)
		m2 := c.addMessage(0)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.Bytecodes[len(c.Bytecodes)-1] == code.Pop {
			c.Bytecodes = c.Bytecodes[:len(c.Bytecodes)-1]
		}

		c.emit(code.Jump)

		m3 := c.addMessage(0)
		m4 := c.addMessage(0)

		c.Messages[m1] = object.NewMessage(len(c.Bytecodes) - 1)
		c.Messages[m2] = object.NewMessage(len(c.Messages))

		if node.Alternative != nil {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.Bytecodes[len(c.Bytecodes)-1] == code.Pop {
				c.Bytecodes = c.Bytecodes[:len(c.Bytecodes)-1]
			}
		} else {
			c.addMessage(c.addConstant(object.NULLOBJ))
			c.emit(code.LoadConstant)
		}

		c.Messages[m3] = object.NewMessage(len(c.Bytecodes) - 1)
		c.Messages[m4] = object.NewMessage(len(c.Messages))

	case *ast.WhileExpression:
		m3i := len(c.Bytecodes) - 1
		m4i := len(c.Messages)

		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		c.emit(code.JumpFalse)
		m1 := c.addMessage(0)
		m2 := c.addMessage(0)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		c.emit(code.Jump)
		m3 := c.addMessage(0)
		m4 := c.addMessage(0)
		c.Messages[m1] = object.NewMessage(len(c.Bytecodes) - 1)
		c.Messages[m2] = object.NewMessage(len(c.Messages))
		c.Messages[m3] = object.NewMessage(m3i)
		c.Messages[m4] = object.NewMessage(m4i)

		c.addMessage(c.addConstant(object.NULLOBJ))
		c.emit(code.LoadConstant)

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
		panic("Node not supported: " + node.TokenLiteral() + " | " + node.String())
	}
	return nil
}
