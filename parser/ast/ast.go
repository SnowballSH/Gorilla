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

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i Integer) e() {}

type Infix struct {
	Left  Node
	Right Node
	Op    token.Token
	Tk    token.Token
}

func (i Infix) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op.Literal, i.Right.String())
}

func (i Infix) e() {}
