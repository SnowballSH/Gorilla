package compiler

import (
	"../ast"
	"../code"
	"../eval"
	"../object"
	"fmt"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) Compile(node ast.Node) error {
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

	case *ast.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			if node.Operator == "<" {
				c.emit(code.Gt)
			} else if node.Operator == "<=" {
				c.emit(code.Gteq)
			}

			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.Add)
		case "-":
			c.emit(code.Sub)
		case "*":
			c.emit(code.Mul)
		case "/":
			c.emit(code.Div)

		case ">":
			c.emit(code.Gt)
		case ">=":
			c.emit(code.Gteq)
		case "==":
			c.emit(code.Eq)
		case "!=":
			c.emit(code.Neq)

		default:
			return fmt.Errorf("[Line %d] unknown operator %s", node.Token.Line, node.Operator)
		}

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.Not)
		case "-":
			c.emit(code.Neg)
		case "+":
			c.emit(code.Pos)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := eval.NewInt(node.Value, node.Token.Line)
		c.emit(code.LoadConst, c.addConstant(integer))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.LoadTrue)
		} else {
			c.emit(code.LoadFalse)
		}

	default:
		return fmt.Errorf("critical error: Not Supported (%T)", node)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
