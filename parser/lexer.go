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
	return &Lexer{
		input:       []rune(input),
		inputLength: utf8.RuneCountInString(input),
		ch:          []rune(input)[0],
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
			ch := l.ch
			l.readChar()
			tok = l.newToken(
				token.Terminator,
				string(ch)+string(l.ch),
			)
		} else {
			tok = l.newToken(token.Terminator, string(l.ch))
		}
		l.linePlace++
		l.charPlace = 0

	case '\n':
		tok = l.newToken(token.Terminator, string(l.ch))
		l.linePlace++
		l.charPlace = 0
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(t string, s string) token.Token {
	return token.Token{
		Type:    t,
		Literal: s,
		Char:    l.charPlace,
		Line:    l.linePlace,
	}
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
