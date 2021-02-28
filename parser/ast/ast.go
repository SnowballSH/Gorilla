package ast

import (
	"fmt"
	"github.com/SnowballSH/Gorilla/parser/token"
)

type Node interface {
	String() string
	Token() token.Token
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
	if e.Es != nil {
		return e.Es.String() + ";"
	}
	return ""
}

func (e ExpressionStatement) Token() token.Token {
	return e.Tk
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

func (i Integer) Token() token.Token {
	return i.Tk
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

func (i Infix) Token() token.Token {
	return i.Tk
}

func (i Infix) e() {}
