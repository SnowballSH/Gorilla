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
func (c *Compiler) emit(b ...byte) {
	c.Result = append(c.Result, b...)
}

// emitString emits a string
func (c *Compiler) emitString(b string) {
	c.emit(byte(len(b)))
	c.emit([]byte(b)...)
}

// emitInt emits an integer
func (c *Compiler) emitInt(b int64) {
	l := leb128.AppendSleb128(nil, b)
	c.emit(byte(len(l)))
	c.emit(l...)
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

	case *ast.Call:
		for _, x := range reverse(e.Arguments) {
			c.compileExpr(x)
		}
		c.compileExpr(e.Function)
		c.emit(grammar.Call)
		c.emitInt(int64(len(e.Arguments)))
	}
}

func reverse(numbers []ast.Expression) []ast.Expression {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
