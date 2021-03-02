package parser

import (
	"github.com/SnowballSH/Gorilla/parser/token"
	"unicode/utf8"
)

type Lexer struct {
	input       []rune
	inputLength int

	ch rune

	position  int
	charPlace int
	linePlace int
}

func NewLexer(input string) *Lexer {
	var ch rune
	if utf8.RuneCountInString(input) == 0 {
		ch = 0
	} else {
		ch = []rune(input)[0]
	}
	return &Lexer{
		input:       []rune(input),
		inputLength: utf8.RuneCountInString(input),
		ch:          ch,
		position:    0,
		charPlace:   1,
		linePlace:   0,
	}
}

func (l *Lexer) next() token.Token {
	l.skipWhitespace()

	var tok token.Token

	switch l.ch {
	case '#':
		for l.peekChar() != '\n' && l.peekChar() != '\r' && l.peekChar() != 0 {
			l.readChar()
		}
		l.readChar()
		return l.next()

	case '\r':
		if l.peekChar() == '\n' {
			l.readChar()
			tok = l.newToken(
				token.Newline,
				"\r\n",
			)
		} else {
			tok = l.newToken(token.Newline, "\n")
		}
		l.linePlace++
		l.charPlace = 0

	case '\n':
		tok = l.newToken(token.Newline, "\n")
		l.linePlace++
		l.charPlace = 0

	case ';':
		tok = l.newToken(token.Semicolon, ";")
	case '=':
		tok = l.newToken(token.Eq, "=")

	case '+':
		tok = l.newToken(token.Plus, "+")
	case '-':
		tok = l.newToken(token.Minus, "-")
	case '*':
		tok = l.newToken(token.Star, "*")
	case '/':
		tok = l.newToken(token.Slash, "/")

	case 0:
		tok = l.newToken(token.EOF, string(byte(0)))

	default:
		if l.isNumber() {
			tok = l.readNumber()
		} else if l.isLetter() || l.ch == '$' {
			tok = l.readIden()
		} else {
			tok = l.newToken(token.Illegal, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() token.Token {
	x := string(l.ch)
	for l.peekIsNumber() {
		l.readChar()
		x += string(l.ch)
	}

	return l.newToken(token.Integer, x)
}

func (l *Lexer) readIden() token.Token {
	x := string(l.ch)
	for l.peekIsLetter() || l.peekIsNumber() {
		l.readChar()
		x += string(l.ch)
	}

	return l.newToken(token.Iden, x)
}

func (l *Lexer) newToken(t string, s string) token.Token {
	return token.Token{
		Type:    t,
		Literal: s,
		Char:    l.charPlace,
		Line:    l.linePlace,
	}
}

func (l *Lexer) isLetter() bool {
	return ('A' <= l.ch && l.ch <= 'Z') || ('a' <= l.ch && l.ch <= 'z') || l.ch == '_'
}

func (l *Lexer) isNumber() bool {
	return '0' <= l.ch && l.ch <= '9'
}

func (l *Lexer) peekIsLetter() bool {
	p := l.peekChar()
	return ('A' <= p && p <= 'Z') || ('a' <= p && p <= 'z') || p == '_'
}

func (l *Lexer) peekIsNumber() bool {
	p := l.peekChar()
	return '0' <= p && p <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' { // || l.ch == '\n' || l.ch == '\r'
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	l.position++
	l.charPlace++
	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
}

func (l *Lexer) peekChar() rune {
	return l.peek(1)
}

func (l *Lexer) peek(amount int) rune {
	if l.position+amount >= len(l.input) {
		return 0
	} else {
		return l.input[l.position+amount]
	}
}
