package compiler

import (
	"Gorilla/ast"
	"Gorilla/code"
	"Gorilla/object"
	"fmt"
)

type MCCombo struct {
	IIndex int
	MIndex int
}

type BytecodeCompiler struct {
	Constants []object.BaseObject
	Bytecodes []code.Opcode
	Messages  []object.Message

	CurrentJumpFalse *MCCombo
	InLoop           bool
	BreakCombos      []*MCCombo
}

func NewBytecodeCompiler() *BytecodeCompiler {
	return &BytecodeCompiler{
		CurrentJumpFalse: nil,
		InLoop:           false,
	}
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
		if len(node.Statements) == 0 {
			c.addMessage(c.addConstant(object.NULLOBJ))
			c.emit(code.LoadConstant)
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

	case *ast.BreakStatement:
		if !c.InLoop {
			return fmt.Errorf("[Line %d] Break outside of loop", node.Token.Line+1)
		}

		c.addMessage(0)
		c.addMessage(0)
		c.emit(code.Jump)
		ii := len(c.Messages) - 2
		mi := len(c.Messages) - 1
		c.BreakCombos = append(c.BreakCombos, &MCCombo{
			IIndex: ii,
			MIndex: mi,
		})

	case *ast.NextStatement:
		if !c.InLoop {
			return fmt.Errorf("[Line %d] Next outside of loop", node.Token.Line+1)
		}

		ii := c.CurrentJumpFalse.IIndex
		mi := c.CurrentJumpFalse.MIndex
		c.addMessage(ii)
		c.addMessage(mi)
		c.emit(code.Jump)

	case *ast.Identifier:
		c.addMessage(node.Value)
		c.addMessage(node.Token.Line)
		c.emit(code.GetVar)

	case *ast.IntegerLiteral:
		c.addMessage(c.addConstant(object.NewInteger(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.FloatLiteral:
		c.addMessage(c.addConstant(object.NewFloat(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.Boolean:
		c.addMessage(c.addConstant(object.NewBool(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.StringLiteral:
		c.addMessage(c.addConstant(object.NewString(node.Value, node.Token.Line)))
		c.emit(code.LoadConstant)

	case *ast.ArrayLiteral:
		for _, v := range node.Elements {
			err := c.Compile(v)
			if err != nil {
				return err
			}
		}
		c.addMessage(len(node.Elements))
		c.addMessage(node.Token.Line)
		c.emit(code.MakeArray)

	case *ast.HashLiteral:
		for k, v := range node.Pairs {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(v)
			if err != nil {
				return err
			}
		}
		c.addMessage(len(node.Pairs))
		c.addMessage(node.Token.Line)
		c.emit(code.MakeHash)

	case *ast.FunctionLiteral:
		e := c.makeFunc(node)
		if e != nil {
			return e
		}

	case *ast.FunctionStmt:
		e := c.makeStmtFunc(node)
		if e != nil {
			return e
		}

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
		case "**":
			name = "pow"
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

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		c.addMessage("getIndex")
		c.emit(code.Method)

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.addMessage(node.Token.Line)
		c.addMessage(1)
		c.emit(code.Call)

	case *ast.IndexAssignmentExpression:
		err := c.Compile(node.Receiver)
		if err != nil {
			return err
		}
		c.addMessage("setIndex")
		c.emit(code.Method)

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		err = c.Compile(node.Value)
		if err != nil {
			return err
		}

		c.addMessage(node.Token.Line)
		c.addMessage(2)
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
		NestLoop := c.InLoop
		LastCombo := c.CurrentJumpFalse
		OldBreakCombo := c.BreakCombos
		c.InLoop = true

		m3i := len(c.Bytecodes) - 1
		m4i := len(c.Messages)

		c.CurrentJumpFalse = &MCCombo{
			IIndex: m3i,
			MIndex: m4i,
		}

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

		for _, combo := range c.BreakCombos {
			c.Messages[combo.IIndex] = c.Messages[m1]
			c.Messages[combo.MIndex] = c.Messages[m2]
		}

		c.CurrentJumpFalse = LastCombo
		c.BreakCombos = OldBreakCombo
		if !NestLoop {
			c.InLoop = false
		}

	case *ast.GetAttr:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}
		c.addMessage(node.Name.String())
		c.emit(code.Method)

	case *ast.AttrAssignmentExpression:
		err := c.Compile(node.Receiver)
		if err != nil {
			return err
		}

		err = c.Compile(node.Value)
		if err != nil {
			return err
		}

		c.emit(code.SetMethod)
		c.addMessage(node.Name)

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

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.Return)

	case *ast.UseStatement:
		c.addMessage(node.Name)
		c.addMessage(node.Token.Line)
		c.emit(code.Import)

	case *ast.DoExpression:
		e := c.makeDo(node)
		if e != nil {
			return e
		}

	default:
		panic("Node not supported: " + node.TokenLiteral() + " | " + node.String())
	}
	return nil
}

func (c *BytecodeCompiler) makeFunc(node *ast.FunctionLiteral) error {
	newc := NewBytecodeCompiler()
	err := newc.Compile(node.Body)
	if err != nil {
		return err
	} else {
		var prms []string
		for _, v := range node.Parameters {
			prms = append(prms, v.Value)
		}

		c.addMessage(c.addConstant(object.NewFunction(
			&object.FunctionValue{
				Constants: newc.Constants,
				Bytecodes: newc.Bytecodes,
				Messages:  newc.Messages,
				Params:    prms,
			}, node.Token.Line)))
		c.emit(code.LoadConstant)
	}
	return nil
}

func (c *BytecodeCompiler) makeDo(node *ast.DoExpression) error {
	newc := NewBytecodeCompiler()
	err := newc.Compile(node.Block)
	if err != nil {
		return err
	} else {
		var prms []string
		for _, v := range node.Params {
			prms = append(prms, v.Value)
		}

		c.addMessage(c.addConstant(object.NewMacro(
			&object.FunctionValue{
				Constants: newc.Constants,
				Bytecodes: newc.Bytecodes,
				Messages:  newc.Messages,
				Params:    prms,
			}, node.Token.Line)))
		c.emit(code.LoadConstant)
	}
	return nil
}

func (c *BytecodeCompiler) makeStmtFunc(node *ast.FunctionStmt) error {
	newc := NewBytecodeCompiler()
	err := newc.Compile(node.Body)
	if err != nil {
		return err
	} else {
		var prms []string
		for _, v := range node.Parameters {
			prms = append(prms, v.Value)
		}

		c.addMessage(c.addConstant(object.NewFunction(
			&object.FunctionValue{
				Constants: newc.Constants,
				Bytecodes: newc.Bytecodes,
				Messages:  newc.Messages,
				Params:    prms,
			}, node.Token.Line)))
		c.emit(code.LoadConstant)
		c.emit(code.SetVar)
		c.addMessage(node.Name)
		c.emit(code.Pop)
	}
	return nil
}
