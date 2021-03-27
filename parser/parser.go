package parser

import (
	"github.com/SnowballSH/Gorilla/errors"
	"github.com/SnowballSH/Gorilla/parser/ast"
	"github.com/SnowballSH/Gorilla/parser/token"
	"strconv"
	"unicode/utf8"
)

// infixPrecedence is the precedence table for Gorilla
var infixPrecedence = map[string][2]byte{
	token.Plus:  {1, 2},
	token.Minus: {1, 2},

	token.Star:  {3, 4},
	token.Slash: {3, 4},
}

// Parser is the base parsing struct
type Parser struct {
	l     *Lexer
	Error *string

	cur  token.Token
	peek token.Token
}

// NewParser creates a parsing from a Lexer
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:     l,
		Error: nil,
	}
	p.next()
	p.next()
	return p
}

// curIs determines whether the current token is the input
func (p *Parser) curIs(t string) bool {
	return p.cur.Type == t
}

// peekIs determines whether the peek token is the input
func (p *Parser) peekIs(t string) bool {
	return p.peek.Type == t
}

// next advances to the next token
func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.l.next()
}

// skipNL skips all Newline
func (p *Parser) skipNL() {
	for p.curIs(token.Newline) {
		p.next()
	}
}

// report panics and reports an error
func (p *Parser) report(why string) {
	err := errors.MakeError(
		string(p.l.input),
		why,
		p.cur.Line, p.cur.Char, utf8.RuneCountInString(p.cur.Literal))
	p.Error = &err
	panic(errors.PARSINGERROR(0))
}

/* ... */

// Parse is the main parsing function. Use this to get an array of Statements.
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

		for p.curIs(token.Semicolon) {
			p.next()
		}

		p.skipNL()
	}
	return program
}

// ParseStatement parses a single statement
func (p *Parser) ParseStatement() ast.Statement {
	p.skipNL()
	switch p.cur.Type {
	default:
		return p.ParseExpressionStatement()
	}
}

// ParseExpressionStatement parses an expression statement
func (p *Parser) ParseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Tk: p.cur}
	stmt.Es = p.ParseExpression(0)
	return stmt
}

// ParseExpression parses an expression
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

// ParseAtom parses all literals
func (p *Parser) ParseAtom() (res ast.Expression) {
	switch p.cur.Type {
	case token.Integer:
		x, e := strconv.ParseInt(p.cur.Literal, 10, 64)
		if e != nil {
			p.report("Could not parse " + processToken(p.cur.Literal) + " as 64-bit integer")
		}

		res = &ast.Integer{
			Value: x,
			Tk:    p.cur,
		}
	case token.String:
		res = &ast.String{
			Value: p.cur.Literal[1 : len(p.cur.Literal)-1],
			Tk:    token.Token{},
		}
	case token.Iden:
		res = p.ParseIden()
	case token.LParen:
		p.next()
		r := p.ParseExpression(0)
		if !p.peekIs(token.RParen) {
			p.next()
			p.report("Expected ')', got " + processToken(p.cur.Literal) + "")
		}
		p.next()
		res = r
	default:
		p.report("Unexpected " + processToken(p.cur.Literal) + "")
		res = nil
	}

	if p.peekIs(token.LParen) {
		var k token.Token

		k = p.peek

		p.next()

		var args []ast.Expression

		for p.peek.Type != token.EOF && p.peek.Type != token.RParen {
			p.next()
			p.skipNL()

			args = append(args, p.ParseExpression(0))

			if p.peekIs(token.RParen) {
				break
			}

			if !p.peekIs(token.Comma) {
				p.next()
				p.report("Expected ',', got " + processToken(p.cur.Literal) + "")
				return
			}

			p.next()
		}

		if p.peekIs(token.RParen) {
			k = p.peek
			p.next()
			p.next()
		} else {
			p.next()
			p.report("Expected ')', got " + processToken(p.cur.Literal) + "")
			return
		}

		res = &ast.Call{
			Function:  res,
			Arguments: args,
			Tk:        k,
		}
	}

	return
}

// ParseIden parses a variable name
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

// processToken is a helper function for errors
func processToken(s string) (ss string) {
	ss = "'" + s + "'"
	if s == "\x00" {
		ss = "End of File"
	}
	return
}
