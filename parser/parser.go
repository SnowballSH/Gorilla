package parser

import (
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/parser/ast"
	"github.com/SnowballSH/Gorilla/parser/token"
	"strconv"
	"unicode/utf8"
)

var infixPrecedence = map[string][2]byte{
	token.Plus:  {1, 2},
	token.Minus: {1, 2},

	token.Star:  {3, 4},
	token.Slash: {3, 4},
}

type Parser struct {
	l     *Lexer
	Error *string

	cur  token.Token
	peek token.Token
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:     l,
		Error: nil,
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
		p.cur.Line, p.cur.Char, utf8.RuneCountInString(p.cur.Literal))
	p.Error = &err
	panic(errors.PARSINGERROR(0))
}

/* ... */

func (p *Parser) Parse() []ast.Statement {
	var program []ast.Statement

	defer func() {
		if r := recover(); r != nil {
			errors.TestERR(r)
		}
	}()

	p.skipNL()
	for !p.curIs(token.EOF) {
		stmt := p.ParseStatement()
		program = append(program, stmt)

		p.next()

		if !p.curIs(token.Newline) && !p.curIs(token.Semicolon) && !p.curIs(token.EOF) {
			p.report("Expected newline or ;, got '" + p.cur.Literal + "'")
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
	stmt.Es = p.ParseExpression(0)
	return stmt
}

func (p *Parser) ParseExpression(pr byte) ast.Expression {
	p.skipNL()
	left := p.ParseAtom()

	for !p.peekIs(token.EOF) {
		op := p.peek
		prs, ok := infixPrecedence[op.Type]
		if !ok {
			break
		}

		if prs[0] < pr {
			break
		}

		p.next()
		p.next()

		right := p.ParseExpression(prs[1])

		left = &ast.Infix{
			Left:  left,
			Right: right,
			Op:    op,
		}
	}

	return left
}

func (p *Parser) ParseAtom() ast.Expression {
	switch p.cur.Type {
	case token.Integer:
		x, e := strconv.ParseInt(p.cur.Literal, 10, 64)
		if e != nil {
			p.report("Could not parse '" + p.cur.Literal + "' as 64-bit integer")
		}

		return &ast.Integer{
			Value: x,
			Tk:    p.cur,
		}
	case token.String:
		return &ast.String{
			Value: p.cur.Literal[1 : len(p.cur.Literal)-1],
			Tk:    token.Token{},
		}
	case token.Iden:
		return p.ParseIden()
	case token.LParen:
		p.next()
		r := p.ParseExpression(0)
		if !p.peekIs(token.RParen) {
			p.report("Expected ')', got '" + p.cur.Literal + "'")
		}
		p.next()
		return r
	default:
		p.report("Unexpected '" + p.cur.Literal + "'")
		return nil
	}
}

func (p *Parser) ParseIden() ast.Expression {
	p.skipNL()
	if p.peekIs(token.Eq) {
		x := p.cur
		p.next()
		p.next()
		return &ast.SetVar{
			Name:  x.Literal,
			Value: p.ParseExpression(0),
			Tk:    x,
		}
	}
	return &ast.GetVar{
		Name: p.cur.Literal,
		Tk:   p.cur,
	}
}
