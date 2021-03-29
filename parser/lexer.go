package parser

import (
	"github.com/SnowballSH/Gorilla/parser/token"
	"unicode/utf8"
)

// Lexer is the tokenizer type
type Lexer struct {
	input       []rune
	inputLength int

	ch rune

	position  int
	charPlace int
	linePlace int
}

// NewLexer creates a lexer from a string
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

// next() advances and returns the next token
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

	case '(':
		tok = l.newToken(token.LParen, "(")
	case ')':
		tok = l.newToken(token.RParen, ")")

	case '{':
		tok = l.newToken(token.LCurly, "{")
	case '}':
		tok = l.newToken(token.RCurly, "}")

	case ',':
		tok = l.newToken(token.Comma, ",")

	case '"':
		x, ok := l.readString('"')
		if !ok {
			tok = l.newToken(token.Illegal, "\"")
		} else {
			tok = l.newToken(token.String, "\""+x+"\"")
		}

	case '\'':
		x, ok := l.readString('\'')
		if !ok {
			tok = l.newToken(token.Illegal, "'")
		} else {
			tok = l.newToken(token.String, "'"+x+"'")
		}

	case 0:
		tok = l.newToken(token.EOF, "\x00")

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

// readNumber() reads a number and returns the number token
func (l *Lexer) readNumber() token.Token {
	x := string(l.ch)
	for l.peekIsNumber() {
		l.readChar()
		x += string(l.ch)
	}

	return l.newToken(token.Integer, x)
}

// readIden() reads an identifier
func (l *Lexer) readIden() token.Token {
	x := string(l.ch)
	for l.peekIsLetter() || l.peekIsNumber() {
		l.readChar()
		x += string(l.ch)
	}

	if k, ok := token.Keywords[x]; ok {
		return l.newToken(k, x)
	}

	return l.newToken(token.Iden, x)
}

// newToken is a helper function to create a new token
func (l *Lexer) newToken(t string, s string) token.Token {
	return token.Token{
		Type:    t,
		Literal: s,
		Char:    l.charPlace,
		Line:    l.linePlace,
	}
}

// isLetter is a helper function to determine whether the current character is a letter or _
func (l *Lexer) isLetter() bool {
	return ('A' <= l.ch && l.ch <= 'Z') || ('a' <= l.ch && l.ch <= 'z') || l.ch == '_'
}

// isNumber is a helper function to determine whether the current character is a number
func (l *Lexer) isNumber() bool {
	return '0' <= l.ch && l.ch <= '9'
}

// peekIsLetter is a helper function to determine whether the peek character is a letter or _
func (l *Lexer) peekIsLetter() bool {
	p := l.peekChar()
	return ('A' <= p && p <= 'Z') || ('a' <= p && p <= 'z') || p == '_'
}

// peekIsNumber is a helper function to determine whether the peek character is a number
func (l *Lexer) peekIsNumber() bool {
	p := l.peekChar()
	return '0' <= p && p <= '9'
}

// skipWhitespace skips newline and tabs and spaces
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' { // || l.ch == '\n' || l.ch == '\r'
		l.readChar()
	}
}

// readChar reads a character and advances
func (l *Lexer) readChar() {
	l.position++
	l.charPlace++
	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}
}

// readString reads a string with escapes
func (l *Lexer) readString(c rune) (string, bool) {
	var res []rune
	for {
		l.readChar()

		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				res = append(res, '\n')
			case 'r':
				res = append(res, '\r')
			case 't':
				res = append(res, '\t')
			case '\\':
				res = append(res, '\\')
			case '"':
				res = append(res, '"')
			case '\'':
				res = append(res, '\'')
			case '`':
				res = append(res, '`')
			case 'v':
				res = append(res, '\v')
			case 'a':
				res = append(res, '\a')
			case 'b':
				res = append(res, '\b')
			default:
				res = append(res, l.ch)
			}
			continue
		}
		if l.ch == '\r' {
			if l.peekChar() == '\n' {
				res = append(res, l.ch)
				l.readChar()
			}
			l.linePlace++
			l.charPlace = 0
		} else if l.ch == '\n' {
			l.linePlace++
			l.charPlace = 0
		}
		if l.ch == c {
			return string(res), true
		}
		if l.ch == 0 {
			return "", false
		}
		res = append(res, l.ch)
	}
}

// peekChar peeks a character
func (l *Lexer) peekChar() rune {
	return l.peek(1)
}

// peek peeks an amount of characters
func (l *Lexer) peek(amount int) rune {
	if l.position+amount >= len(l.input) {
		return 0
	} else {
		return l.input[l.position+amount]
	}
}
