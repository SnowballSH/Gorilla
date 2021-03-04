package ast

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/parser/token"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
	s()
}

type Expression interface {
	Node
	e()
}

/* ... */

type ExpressionStatement struct {
	Es Expression
	Tk token.Token
}

func (e ExpressionStatement) String() string {
	return e.Es.String() + ";"
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

func (i *Integer) e() {}

type GetVar struct {
	Name string
	Tk   token.Token
}

func (g *GetVar) String() string {
	return g.Name
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

func (*SetVar) e() {}

type Infix struct {
	Left  Node
	Right Node
	Op    token.Token
}

func (i *Infix) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op.Literal, i.Right.String())
}

func (i *Infix) e() {}
