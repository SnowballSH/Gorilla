package parser

import (
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/parser/ast"
	"github.com/SnowballSH/Gorilla/parser/token"
	"strconv"
)

var infixPrecedence = map[string][2]byte{
	token.Plus:  {1, 2},
	token.Minus: {1, 2},

	token.Star:  {3, 4},
	token.Slash: {3, 4},
}

type Parser struct {
	l     *Lexer
	error *string

	cur  token.Token
	peek token.Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:     l,
		error: nil,
	}
	p.next()
	p.next()
	return p
}

func (p *Parser) curIs(t string) bool {
	return p.cur.Type == t
}

func (p *Parser) peekIs(t string) bool {
	return p.peek.Type == t
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.next()
}

func (p *Parser) skipNL() {
	for p.curIs(token.Newline) {
		p.next()
	}
}

func (p *Parser) report(why string) {
	err := errors.MakeError(
		string(p.l.input),
		why,
		p.cur.Line, p.cur.Char, len(p.cur.Literal))
	p.error = &err
}

/* ... */

func (p *Parser) Parse() []ast.Node {
	var program []ast.Node

	p.skipNL()
	for !p.curIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			program = append(program, stmt)
		} else {
			return nil
		}

		p.next()

		if !p.curIs(token.Newline) && !p.curIs(token.Semicolon) && !p.curIs(token.EOF) {
			p.report("Expected newline or ;, got " + p.cur.Literal)
			return nil
		}

		p.skipNL()
	}
	return program
}

func (p *Parser) ParseStatement() ast.Statement {
	p.skipNL()
	switch p.cur.Type {
	default:
		return p.ParseExpressionStatement()
	}
}

func (p *Parser) ParseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Tk: p.cur}
	stmt.Es = p.ParseExpression()
	if stmt.Es == nil {
		return nil
	}
	return stmt
}

func (p *Parser) ParseExpression() ast.Expression {
	switch p.cur.Type {
	case token.Integer:
		x, e := strconv.ParseInt(p.cur.Literal, 10, 64)
		if e != nil {
			p.report("Could not parse '" + p.cur.Literal + "' as 64-bit integer")
			return nil
		}

		return &ast.Integer{
			Value: x,
			Tk:    p.cur,
		}
	default:
		panic("error")
	}
}

func (p *Parser) parseInfixExpr() ast.Node {
	panic(0)
}
