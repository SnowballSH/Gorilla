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
	token.Star:    {25, 26},
	token.Slash:   {25, 26},
	token.Percent: {25, 26},

	token.Plus:  {23, 24},
	token.Minus: {23, 24},

	token.Smaller:   {19, 20},
	token.Larger:    {19, 20},
	token.SmallerEq: {19, 20},
	token.LargerEq:  {19, 20},

	token.DbEq: {17, 18},
	token.Neq:  {17, 18},
}

// prefixPrecedence is the prefix precedence table for Gorilla
var prefixPrecedence = map[string]byte{
	token.Plus:  27,
	token.Minus: 27,
	token.Not:   27,
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

// skipNLPeek skips all Newline in peek
func (p *Parser) skipNLPeek() {
	for p.peekIs(token.Newline) {
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
	var w ast.Expression

	w = p.ParseExpression(0)

	stmt.Es = w
	return stmt
}

// ParseBlock parses a block
func (p *Parser) ParseBlock() *ast.Block {
	if !p.curIs(token.LCurly) {
		x := p.ParseStatement()
		return &ast.Block{Stmts: []ast.Statement{x}}
	}

	var stmts = &ast.Block{Stmts: nil}

	p.next()

	p.skipNL()

	for !p.curIs(token.RCurly) && !p.curIs(token.EOF) {
		stmt := p.ParseStatement()
		stmts.Stmts = append(stmts.Stmts, stmt)

		p.next()

		if !p.curIs(token.Newline) && !p.curIs(token.Semicolon) && !p.curIs(token.EOF) && !p.curIs(token.RCurly) {
			p.report("Expected newline or ;, got '" + p.cur.Literal + "'")
		}

		p.skipNL()

		for p.curIs(token.Semicolon) {
			p.next()
		}

		p.skipNL()
	}

	if !p.curIs(token.RCurly) {
		p.report("Expected '}', got " + processToken(p.cur.Literal))
	}

	return stmts
}

// ParseIfElse parses an if, else expression
func (p *Parser) ParseIfElse() *ast.IfElse {
	tk := p.cur
	p.next()
	cond := p.ParseExpression(0)
	p.next()
	If := p.ParseBlock()

	if p.peekIs(token.Else) {
		p.next()
		p.next()
		Else := p.ParseBlock()
		return &ast.IfElse{
			Condition: cond,
			If:        If,
			Else:      Else,
			Tk:        tk,
		}
	}

	return &ast.IfElse{
		Condition: cond,
		If:        If,
		Else:      &ast.Block{Stmts: nil},
		Tk:        tk,
	}
}

// ParseIfElse parses a lambda expression
// |a, b, c| { do_something(a * b * c) }
func (p *Parser) ParseLambda() *ast.Lambda {
	tk := p.cur // |

	var idens []string

	for p.peek.Type != token.EOF && p.peek.Type != token.VBar {
		p.next()
		p.skipNL()

		iden := p.ParseIdenString()
		idens = append(idens, iden)
		p.skipNLPeek()

		if p.peekIs(token.VBar) {
			break
		}

		if !p.peekIs(token.Comma) {
			p.next()
			p.report("Expected ',', got " + processToken(p.cur.Literal))
		}

		p.next()

		p.skipNLPeek()
	}

	p.next()
	p.skipNL()
	if p.curIs(token.VBar) {
	} else {
		p.report("Expected '|', got " + processToken(p.cur.Literal))
	}

	p.next()
	block := p.ParseBlock()

	return &ast.Lambda{
		Arguments: idens,
		Block:     block,
		Tk:        tk,
	}
}

// ParseExpression parses an expression
func (p *Parser) ParseExpression(pr byte) ast.Expression {
	p.skipNL()
	var left ast.Expression
	if prs, ok := prefixPrecedence[p.cur.Type]; ok {
		tk := p.cur
		p.next()
		left = &ast.Prefix{
			Right: p.ParseExpression(prs),
			Op:    tk,
		}
	} else {
		switch p.cur.Type {
		case token.If:
			left = p.ParseIfElse()
		case token.VBar:
			left = p.ParseLambda()
		case token.LCurly:
			k := p.cur
			left = &ast.Closure{
				Block: p.ParseBlock(),
				Tk:    k,
			}
		default:
			left = p.ParseAtom()
		}

		// suffix
		for {
			if p.peekIs(token.LParen) {
				var k token.Token

				k = p.peek

				p.next()

				var args []ast.Expression

				for p.peek.Type != token.EOF && p.peek.Type != token.RParen {
					p.next()
					p.skipNL()

					args = append(args, p.ParseExpression(0))

					p.skipNLPeek()

					if p.peekIs(token.RParen) {
						break
					}

					if !p.peekIs(token.Comma) {
						p.next()
						p.report("Expected ',', got " + processToken(p.cur.Literal) + "")
					}

					p.next()

					p.skipNLPeek()
				}

				p.next()
				p.skipNL()
				if p.curIs(token.RParen) {
					k = p.cur
				} else {
					p.report("Expected ')', got " + processToken(p.cur.Literal) + "")
				}

				left = &ast.Call{
					Function:  left,
					Arguments: args,
					Tk:        k,
				}
			} else if p.peekIs(token.Dot) {
				k := p.peek
				p.next()
				p.next()
				str := p.ParseIdenString()

				left = &ast.GetInstance{
					Parent: left,
					Name:   str,
					Tk:     k,
				}
			} else {
				break
			}
		}
	}

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
			p.report("Expected ')', got " + processToken(p.cur.Literal))
		}
		p.next()
		res = r
	default:
		p.report("Unexpected " + processToken(p.cur.Literal))
		res = nil
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

func (p *Parser) ParseIdenString() string {
	p.skipNL()
	x := p.cur
	if x.Type != token.Iden {
		p.report("Expected Identifier, got " + processToken(x.Type))
	}

	return x.Literal
}

// processToken is a helper function for errors
func processToken(s string) (ss string) {
	ss = "'" + s + "'"
	if s == "\x00" {
		ss = "End of File"
	}
	return
}
