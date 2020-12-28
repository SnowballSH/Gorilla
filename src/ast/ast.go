package ast

import (
	"bytes"
	"strings"

	"../token"
)

// Node defines an interface for all nodes in the AST.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement defines the interface for all statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression defines the interface for all expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node. All programs consist of a slice of Statement(s)
type Program struct {
	Statements []Statement
}

// TokenLiteral prints the literal value of the token associated with this node
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String returns a stringified version of the AST for debugging
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// LetStatement the `let` statement represents the AST node that binds an
// expression to an identifier
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ReturnStatement represenets the `return` statement node
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// String returns a stringified version of the AST for debugging
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement represents an expression statement and holds an
// expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// String returns a stringified version of the AST for debugging
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a block statement and holds one or more other
// statements
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }

// String returns a stringified version of the AST for debugging
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier represents an identiifer and holds the name of the identifier
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// String returns a stringified version of the AST for debugging
func (i *Identifier) String() string { return i.Value }

// Boolean represents a boolean value and holds the underlying boolean value
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }

// String returns a stringified version of the AST for debugging
func (b *Boolean) String() string { return b.Token.Literal }

// IntegerLiteral represents a literal integare and holds an integer value
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the AST for debugging
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// PrefixExpression represents a prefix expression and holds the operator
// as well as the right-hand side expression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns a stringified version of the AST for debugging
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// InfixExpression represents an infix expression and holds the left-hand
// expression, operator and right-hand expression
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// IfExpression represents an `if` expression and holds the condition,
// consequence and alternative expressions
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// FunctionLiteral represents a literal functions and holds the function's
// formal parameters and boy of the function as a block statement
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a stringified version of the AST for debugging
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents a call expression and holds the function to be
// called as well as the arguments to be passed to that function
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
