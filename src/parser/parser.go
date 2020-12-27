package parser

import (
	"fmt"
	//"strconv"

	"../ast"
	"../lexer"
	"../token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens
	p.nextToken()
	p.nextToken()

	return p
}

// HELPER START

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	println(p.curToken.Line)
}

func (p *Parser) curTokenIs(t token.TType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) expectPeek(t token.TType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}
func (p *Parser) peekError(t token.TType) {
	msg := fmt.Sprintf("[Line %d] expected %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// HELPER END

// RULES START

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		for p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}
