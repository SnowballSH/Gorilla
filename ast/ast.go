package ast

import (
	"bytes"
	"strings"

	"Gorilla/token"
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

// ReturnStatement represents the `return` statement node
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

type UseStatement struct {
	Token token.Token // the 'use' token
	Name  string
}

func (us *UseStatement) statementNode() {}

func (us *UseStatement) TokenLiteral() string { return us.Token.Literal }

func (us *UseStatement) String() string {
	var out bytes.Buffer

	out.WriteString(us.TokenLiteral() + " ")

	out.WriteString(us.Name)

	out.WriteString(";")

	return out.String()
}

type BreakStatement struct {
	Token token.Token
}

func (s *BreakStatement) statementNode() {}

func (s *BreakStatement) TokenLiteral() string { return s.Token.Literal }

func (s *BreakStatement) String() string {
	return "break;"
}

type NextStatement struct {
	Token token.Token
}

func (s *NextStatement) statementNode() {}

func (s *NextStatement) TokenLiteral() string { return s.Token.Literal }

func (s *NextStatement) String() string {
	return "next;"
}

// FunctionStmt
type FunctionStmt struct {
	Token      token.Token
	Name       string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fs *FunctionStmt) statementNode() {}

func (fs *FunctionStmt) TokenLiteral() string { return fs.Token.Literal }

func (fs *FunctionStmt) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral() + " ")
	out.WriteString(fs.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fs.Body.String())

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
		return es.Expression.String() + ";"
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

// Identifier represents an identifier and holds the name of the identifier
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

// IntegerLiteral represents a literal integer and holds an integer value
type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

// String returns a stringified version of the AST for debugging
func (il *IntegerLiteral) String() string { return il.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

// String returns a stringified version of the AST for debugging
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

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

// IfExpression
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

// WhileExpression
type WhileExpression struct {
	Token       token.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode() {}

// TokenLiteral
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }

// String returns a stringified version of the AST for debugging
func (we *WhileExpression) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())

	return out.String()
}

// AssignmentExpression
type AssignmentExpression struct {
	Token token.Token // The ident token
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}

// TokenLiteral
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Name.String())
	out.WriteString(" = ")

	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	} else {
		out.WriteString("null")
	}

	return out.String()
}

// IndexAssignmentExpression
type IndexAssignmentExpression struct {
	Token    token.Token // The ident token
	Receiver Expression
	Index    Expression
	Value    Expression
}

func (ae *IndexAssignmentExpression) expressionNode() {}

// TokenLiteral
func (ae *IndexAssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ae *IndexAssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Receiver.String())
	out.WriteString("[")
	out.WriteString(ae.Index.String())
	out.WriteString("] = ")

	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	} else {
		out.WriteString("null")
	}

	return out.String()
}

// AttrAssignmentExpression
type AttrAssignmentExpression struct {
	Token    token.Token // The ident token
	Receiver Expression
	Name     string
	Value    Expression
}

func (ae *AttrAssignmentExpression) expressionNode() {}

// TokenLiteral
func (ae *AttrAssignmentExpression) TokenLiteral() string { return ae.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ae *AttrAssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ae.Receiver.String())
	out.WriteString(".")
	out.WriteString(ae.Name)
	out.WriteString(" = ")

	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	} else {
		out.WriteString("null")
	}

	return out.String()
}

// FunctionLiteral
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

// GetAttr is x.y
type GetAttr struct {
	Token token.Token
	Expr  Expression
	Name  Expression
}

func (ge *GetAttr) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ge *GetAttr) TokenLiteral() string { return ge.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ge *GetAttr) String() string {
	var out bytes.Buffer

	out.WriteString(ge.Expr.String())
	out.WriteString(".")
	out.WriteString(ge.Name.String())

	return out.String()
}

// StringLiteral is the string ast
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "`" + sl.Token.Literal + "`" }

// ArrayLiteral represents the array literal and holds a list of expressions
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }

// String returns a stringified version of the AST for debugging
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
