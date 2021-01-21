package compiler

import (
	"../ast"
	"../code"
	"../object"
	"fmt"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Scope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants []object.Object

	scopes     []Scope
	scopeIndex int

	symbolTable *SymbolTable
}

func New() *Compiler {
	mainScope := Scope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants: []object.Object{},

		scopes:     []Scope{mainScope},
		scopeIndex: 0,

		symbolTable: symbolTable,
	}
}

func NewWithState(symbolTable *SymbolTable, constants []object.Object) *Compiler {
	c := New()
	c.symbolTable = symbolTable
	c.constants = constants
	return c
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.LoadGlobal, s.Index)
	case LocalScope:
		c.emit(code.LoadLocal, s.Index)
	case BuiltinScope:
		c.emit(code.LoadBuiltin, s.Index)
	case FreeScope:
		c.emit(code.LoadFree, s.Index)
	}
}

func (c *Compiler) enterScope() {
	scope := Scope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	n := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = n
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.Ret))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.Ret
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

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

	case *ast.LetStatement:
		var (
			ok     bool
			symbol Symbol
		)

		symbol, ok = c.symbolTable.Resolve(node.Name.Value)
		if !ok {
			symbol = c.symbolTable.Define(node.Name.Value)
		}

		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		if symbol.Scope == GlobalScope {
			c.emit(code.SetGlobal, symbol.Index)
		} else {
			c.emit(code.SetLocal, symbol.Index)
		}

	case *ast.FunctionStmt:
		var (
			ok     bool
			symbol Symbol
		)

		symbol, ok = c.symbolTable.Resolve(node.Name)
		if !ok {
			symbol = c.symbolTable.Define(node.Name)
		}

		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.Pop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.Ret) {
			c.emit(code.RetNull)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(code.Closure, fnIndex, len(freeSymbols))

		if symbol.Scope == GlobalScope {
			c.emit(code.SetGlobal, symbol.Index)
		} else {
			c.emit(code.SetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("[Line %d] Variable '%s' is not defined", node.Token.Line+1, node.Value)
		}

		c.loadSymbol(symbol)

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

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `JumpIfFalse` with a bogus value
		jumpIfFalsePos := c.emit(code.JumpElse, 0xFFFF)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.Pop) {
			c.removeLastPop()
		}

		// Emit an `Jump` with a bogus value
		jumpPos := c.emit(code.Jump, 0xFFFF)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpIfFalsePos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.LoadNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.Pop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

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

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.Gt)
			return nil
		}

		if node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.Gteq)
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

		case "<-":
			c.emit(code.LARR)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.StringLiteral:
		str := object.NewString(node.Value, node.Token.Line)
		c.emit(code.LoadConst, c.addConstant(str))

	case *ast.IntegerLiteral:
		integer := object.NewInt(node.Value, node.Token.Line)
		c.emit(code.LoadConst, c.addConstant(integer))

	case *ast.FunctionLiteral:
		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.Pop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.Ret) {
			c.emit(code.RetNull)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(code.Closure, fnIndex, len(freeSymbols))

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(code.Call, len(node.Arguments))

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.Ret)

	case *ast.Boolean:
		if node.Value {
			c.emit(code.LoadTrue)
		} else {
			c.emit(code.LoadFalse)
		}

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.Array, len(node.Elements), node.Token.Line)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.Index)

	case *ast.GetAttr:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}

		c.emit(code.LoadConst, c.addConstant(object.NewString(node.Name.String(), node.Token.Line)))

		c.emit(code.GetAttr)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
