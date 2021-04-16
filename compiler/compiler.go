package compiler

import (
	"ekyu.moe/leb128"
	"github.com/SnowballSH/Gorilla/grammar"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

// Compiler is the base compiler struct
type Compiler struct {
	Result   []byte
	lastLine int
}

// NewCompiler creates a compiler with Magic
func NewCompiler() *Compiler {
	return &Compiler{
		Result:   []byte{grammar.Magic},
		lastLine: 0,
	}
}

// updateLine updates the line count using grammar.Advance and grammar.Back
func (c *Compiler) updateLine(line int) {
	for line > c.lastLine {
		c.lastLine++
		c.emit(grammar.Advance)
	}
	for line < c.lastLine {
		c.lastLine--
		c.emit(grammar.Back)
	}
}

// emit emits some bytes
func (c *Compiler) emit(b ...byte) (pos int) {
	pos = len(c.Result)
	c.Result = append(c.Result, b...)
	return
}

// emitString emits a string
func (c *Compiler) emitString(b string) int {
	c.emit(byte(len(b)))
	c.emit([]byte(b)...)

	return len(b) + 1
}

// emitInt emits an integer
func (c *Compiler) emitInt(b int64) int {
	l := leb128.AppendSleb128(nil, b)
	c.emit(byte(len(l)))
	c.emit(l...)

	return len(l) + 1
}

// Compiler is the base compiling function.
// This function compiles nodes into Result
func (c *Compiler) Compile(nodes []ast.Statement) {
	for _, x := range nodes {
		c.compileNode(x)
	}
}

// compileNode compiles a single node
func (c *Compiler) compileNode(node ast.Node) {
	c.updateLine(node.Line())

	switch v := node.(type) {
	case *ast.ExpressionStatement:
		c.compileExpr(v.Es)
		c.emit(grammar.Pop)
	}
}

// compilerExpr compiles an expression
func (c *Compiler) compileExpr(v ast.Expression) {
	c.updateLine(v.Line())

	switch e := v.(type) {
	case *ast.Integer:
		c.emit(grammar.Integer)

		c.emitInt(e.Value)

	case *ast.String:
		c.emit(grammar.String)
		c.emitString(e.Value)

	case *ast.GetVar:
		c.emit(grammar.GetVar)

		c.emitString(e.Name)

	case *ast.SetVar:
		c.compileExpr(e.Value)

		c.updateLine(e.Line())

		c.emit(grammar.SetVar)

		c.emitString(e.Name)

	case *ast.Infix:
		c.compileExpr(e.Right)
		c.updateLine(e.Line())

		c.compileExpr(e.Left)
		c.updateLine(e.Line())

		c.emit(grammar.GetInstance)
		c.emitString(e.Op.Literal)

		c.emit(grammar.Call)
		c.emitInt(1)

	case *ast.Prefix:
		c.compileExpr(e.Right)
		c.updateLine(e.Line())

		c.emit(grammar.GetInstance)
		c.emitString(e.Op.Literal + "@")

		c.emit(grammar.Call)
		c.emitInt(0)

	case *ast.Call:
		for _, x := range reverse(e.Arguments) {
			c.compileExpr(x)
		}
		c.compileExpr(e.Function)
		c.emit(grammar.Call)
		c.emitInt(int64(len(e.Arguments)))

	case *ast.IfElse:
		c.compileExpr(e.Condition)

		compi := &Compiler{
			Result:   nil,
			lastLine: c.lastLine,
		}
		compe := &Compiler{
			Result:   nil,
			lastLine: c.lastLine,
		}

		ci := make(chan bool)
		ce := make(chan bool)

		go func(ch chan bool) {
			compi.Compile(e.If.Stmts)
			ch <- true
		}(ci)

		go func(ch chan bool) {
			compe.Compile(e.Else.Stmts)
			ch <- true
		}(ce)

		<-ci
		<-ce

		if len(compi.Result) > 0 && compi.Result[len(compi.Result)-1] == grammar.Pop {
			compi.Result[len(compi.Result)-1] = grammar.Noop
		} else {
			compi.emit(grammar.Null)
		}

		if len(compe.Result) > 0 && compe.Result[len(compe.Result)-1] == grammar.Pop {
			compe.Result[len(compe.Result)-1] = grammar.Noop
		} else {
			compe.emit(grammar.Null)
		}

		compi.updateLine(v.Line())
		compe.updateLine(v.Line())

		c.emit(grammar.JumpIfFalse)

		compi.emit(grammar.Jump)
		compi.emitInt(int64(len(c.Result) + 2 + len(leb128.AppendSleb128(nil, int64(len(c.Result)+len(compi.Result)+1))) +
			len(compi.Result) + len(compe.Result)))

		c.emitInt(int64(len(c.Result) + len(compi.Result) + 1))

		c.emit(compi.Result...)
		c.emit(compe.Result...)

	case *ast.Lambda:
		c.emit(grammar.Lambda)
		c.emitInt(int64(len(e.Arguments)))
		for _, x := range e.Arguments {
			c.emitString(x)
		}

		comp := &Compiler{
			Result:   []byte{grammar.Magic},
			lastLine: c.lastLine,
		}
		comp.Compile(e.Block.Stmts)

		c.emitInt(int64(len(comp.Result)))

		c.emit(comp.Result...)

		c.updateLine(comp.lastLine)

	case *ast.GetInstance:
		c.compileExpr(e.Parent)
		c.emit(grammar.GetInstance)
		c.emitString(e.Name)
	}
}

func reverse(numbers []ast.Expression) []ast.Expression {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
