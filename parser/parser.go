package parser

import (
	"fmt"
	"strconv"

	"github.com/SnowballSH/Gorilla/ast"
	"github.com/SnowballSH/Gorilla/lexer"
	"github.com/SnowballSH/Gorilla/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ARR // <-
	OR
	AND
	EQUALS      // ==
	LESSGREATER // > or < or >= or <=
	SUM         // +
	PRODUCT     // *
	POWER       // **
	PREFIX      // -X or !X
	DOT         // a.b
	CALL        // myFunction(X)
	INDEX       // x[y]
)

var precedences = map[token.TType]int{
	token.LARR:     ARR,
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTEQ:     LESSGREATER,
	token.GTEQ:     LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MOD:      PRODUCT,

	token.POW: POWER,

	token.OR:  OR,
	token.AND: AND,

	token.LPAREN:   CALL,
	token.DO:       CALL,
	token.FUNCTION: CALL,
	token.DOT:      DOT,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TType]prefixParseFn
	infixParseFns  map[token.TType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.PLUS, p.parsePrefixExpression)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.DO, p.parseDoExpression)

	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.POW, p.parseInfixExpression)

	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTEQ, p.parseInfixExpression)
	p.registerInfix(token.GTEQ, p.parseInfixExpression)

	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.DO, p.parseCallOneExpression)
	p.registerInfix(token.FUNCTION, p.parseCallOneExpression)
	p.registerInfix(token.DOT, p.parseGetAttr)

	p.registerInfix(token.LARR, p.parseInfixExpression)

	// Read two tokens
	p.nextToken()
	p.nextToken()

	return p
}

// HELPER START

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
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
		p.peekToken.Line+1, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curError(t token.TType) {
	msg := fmt.Sprintf("[Line %d] expected %s, got %s instead",
		p.curToken.Line+1, t, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t token.TType) {
	msg := fmt.Sprintf("[Line %d] Invalid Syntax: unexpected '%s'", p.curToken.Line+1, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// HELPER END

// RULES START

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()

		if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
			msg := fmt.Sprintf(
				"[Line %d] Invalid Syntax: Expected Newline or ';', got '%s'",
				p.curToken.Line+1, p.curToken.Type,
			)
			p.errors = append(p.errors, msg)
		}

		for p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.STMTFUNCTION:
		return p.parseFunctionStatement()
	case token.BREAK:
		return p.parseBreak()
	case token.NEXT:
		return p.parseNext()
	case token.USE:
		return p.parseUseStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.parseIdentifierIden()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseBreak() *ast.BreakStatement {
	stmt := &ast.BreakStatement{Token: p.curToken}

	p.nextToken()
	return stmt
}

func (p *Parser) parseNext() *ast.NextStatement {
	stmt := &ast.NextStatement{Token: p.curToken}

	p.nextToken()
	return stmt
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	p.nextToken()
	res := p.parseBlockStatement()
	if res != nil {
		expression.Consequence = res
	} else {
		return nil
	}

	return expression
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseUseStatement() *ast.UseStatement {
	stmt := &ast.UseStatement{Token: p.curToken}

	if !p.expectPeek(token.STRING) {
		return nil
	}

	iden := p.parseStringLiteral()

	stmt.Name = iden.TokenLiteral()

	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStmt {
	lit := &ast.FunctionStmt{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	lit.Name = p.curToken.Literal

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		lit.Parameters = p.parseFunctionParameters()
	} else {
		lit.Parameters = []*ast.Identifier{}
	}

	p.nextToken()

	res := p.parseBlockStatement()
	if res != nil {
		lit.Body = res
	} else {
		return nil
	}

	return lit
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.peekTokenIs(token.ASSIGN) {
		stmt := &ast.AssignmentExpression{Token: p.curToken}

		stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		p.nextToken()

		p.nextToken()

		stmt.Value = p.parseExpression(LOWEST)

		return stmt
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIdentifierIden() *ast.Identifier {
	if p.curToken.Type != token.IDENT {
		return nil
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("[Line %d] Could not parse %q as integer", p.curToken.Line, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("[Line %d] Could not parse %q as float", p.curToken.Line, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	p.nextToken()
	res := p.parseBlockStatement()
	if res != nil {
		expression.Consequence = res
	} else {
		return nil
	}

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		p.nextToken()
		res = p.parseBlockStatement()
		if res != nil {
			expression.Alternative = res
		} else {
			return nil
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if !p.curTokenIs(token.LBRACE) {
		res := p.parseStatement()
		if res == nil {
			return nil
		}
		block.Statements = append(block.Statements, res)
		return block
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
		for p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.EOF) {
		p.curError(token.RBRACE)
		return nil
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		lit.Parameters = p.parseFunctionParameters()
	} else {
		lit.Parameters = []*ast.Identifier{}
	}

	p.nextToken()
	res := p.parseBlockStatement()
	if res != nil {
		lit.Body = res
	} else {
		return nil
	}

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := p.parseIdentifierIden()
	if ident == nil {
		return nil
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident = p.parseIdentifierIden()
		if ident == nil {
			return nil
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallOneExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = []ast.Expression{p.parseExpression(LOWEST)}
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	return p.parseList(token.RPAREN)
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseList(token.RBRACKET)

	return array
}

func (p *Parser) parseList(t token.TType) []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(t) {
		p.nextToken()
		return args
	}

	p.nextToken()

	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		for p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		args = append(args, p.parseExpression(LOWEST))
		for p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if !p.expectPeek(t) {
		return nil
	}
	return args
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		return &ast.IndexAssignmentExpression{Token: exp.Token, Receiver: left, Index: exp.Index, Value: value}
	}

	return exp
}

func (p *Parser) parseGetAttr(expr ast.Expression) ast.Expression {
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	crr := p.curToken
	iden := p.parseIdentifierIden()

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		value := p.parseExpression(LOWEST)
		return &ast.AttrAssignmentExpression{
			Token:    crr,
			Receiver: expr,
			Name:     iden.String(),
			Value:    value,
		}
	}

	return &ast.GetAttr{
		Token: crr,
		Expr:  expr,
		Name:  iden,
	}
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}

	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseDoExpression() ast.Expression {
	lit := &ast.DoExpression{Token: p.curToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		lit.Params = p.parseFunctionParameters()
	} else {
		lit.Params = []*ast.Identifier{}
	}

	p.nextToken()

	res := p.parseBlockStatement()
	if res != nil {
		lit.Block = res
	} else {
		return nil
	}

	return lit
}
