package compiler

import (
	"ekyu.moe/leb128"
	"github.com/SnowballSH/Gorilla/grammar"

	"github.com/SnowballSH/Gorilla/parser/ast"
)

type Compiler struct {
	Result   []byte
	lastLine int
}

func NewCompiler() *Compiler {
	return &Compiler{
		Result:   []byte{grammar.Magic},
		lastLine: 0,
	}
}

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

func (c *Compiler) emit(b ...byte) {
	c.Result = append(c.Result, b...)
}

func (c *Compiler) emitString(b string) {
	c.emit(byte(len(b)))
	c.emit([]byte(b)...)
}

func (c *Compiler) emitInt(b int64) {
	l := leb128.AppendSleb128(nil, b)
	c.emit(byte(len(l)))
	c.emit(l...)
}

func (c *Compiler) Compile(nodes []ast.Statement) {
	for _, x := range nodes {
		c.compileNode(x)
	}
}

func (c *Compiler) compileNode(node ast.Node) {
	c.updateLine(node.Line())

	switch v := node.(type) {
	case *ast.ExpressionStatement:
		c.compileExpr(v.Es)
		c.emit(grammar.Pop)
	}
}

func (c *Compiler) compileExpr(v ast.Expression) {
	c.updateLine(v.Line())

	switch e := v.(type) {
	case *ast.Integer:
		c.emit(grammar.Integer)

		c.emitInt(e.Value)

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
	}
}
