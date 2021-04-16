package ast

import (
	"bytes"
	"fmt"
	"github.com/SnowballSH/Gorilla/parser/token"
	"strings"
)

// All ast nodes inherits from Node
type Node interface {
	String() string
	Line() int
}

// Statement represents a statement
type Statement interface {
	Node
	s()
}

// Expression represents an expression
type Expression interface {
	Node
	e()
}

type Block struct {
	Stmts []Statement
}

func (b *Block) String() string {
	x := "{"
	for _, y := range b.Stmts {
		x += "\n" + y.String()
	}
	return x + "\n}"
}

/* ... */

type ExpressionStatement struct {
	Es Expression
	Tk token.Token
}

func (e ExpressionStatement) String() string {
	return e.Es.String() + ";"
}

func (e ExpressionStatement) Line() int {
	return e.Tk.Line
}

func (e ExpressionStatement) s() {}

/* ... */

type Integer struct {
	Value int64
	Tk    token.Token
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Line() int {
	return i.Tk.Line
}

func (i *Integer) e() {}

type String struct {
	Value string
	Tk    token.Token
}

func (s *String) String() string {
	return fmt.Sprintf("'%s'", s.Value)
}

func (s *String) Line() int {
	return s.Tk.Line
}

func (s *String) e() {}

type GetVar struct {
	Name string
	Tk   token.Token
}

func (g *GetVar) String() string {
	return g.Name
}

func (g *GetVar) Line() int {
	return g.Tk.Line
}

func (*GetVar) e() {}

type SetVar struct {
	Name  string
	Value Expression
	Tk    token.Token
}

func (s *SetVar) String() string {
	return fmt.Sprintf("(%s = %s)", s.Name, s.Value)
}

func (s *SetVar) Line() int {
	return s.Tk.Line
}

func (*SetVar) e() {}

type Infix struct {
	Left  Expression
	Right Expression
	Op    token.Token
}

func (i *Infix) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op.Literal, i.Right.String())
}

func (i *Infix) e() {}

func (i *Infix) Line() int {
	return i.Op.Line
}

type Prefix struct {
	Right Expression
	Op    token.Token
}

func (p *Prefix) String() string {
	return fmt.Sprintf("(%s%s)", p.Op.Literal, p.Right.String())
}

func (p *Prefix) e() {}

func (p *Prefix) Line() int {
	return p.Op.Line
}

type Call struct {
	Function  Expression
	Arguments []Expression
	Tk        token.Token
}

func (c *Call) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return "(" + out.String() + ")"
}

func (c *Call) e() {}

func (c *Call) Line() int {
	return c.Tk.Line
}

type IfElse struct {
	Condition Expression
	If        *Block
	Else      *Block
	Tk        token.Token
}

func (ie *IfElse) String() string {
	return fmt.Sprintf("(if %s %s else %s)", ie.Condition.String(), ie.If.String(), ie.Else.String())
}

func (ie *IfElse) e() {}

func (ie *IfElse) Line() int {
	return ie.Tk.Line
}

type Lambda struct {
	Arguments []string
	Block     *Block
	Tk        token.Token
}

func (l *Lambda) String() string {
	return fmt.Sprintf("(|%s| %s)", strings.Join(l.Arguments, ", "), l.Block.String())
}

func (l *Lambda) e() {}

func (l *Lambda) Line() int {
	return l.Tk.Line
}

type GetInstance struct {
	Parent Expression
	Name   string
	Tk     token.Token
}

func (g *GetInstance) String() string {
	return fmt.Sprintf("(%s.%s)", g.Parent.String(), g.Name)
}

func (g *GetInstance) e() {}

func (g *GetInstance) Line() int {
	return g.Tk.Line
}
