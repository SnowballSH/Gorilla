package lexer

import "../token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	lineCount    int  // current # of line
}

func New(input string) *Lexer {
	l := &Lexer{input: input, lineCount: 0}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	tok.Line = l.lineCount

	l.skipWhitespace()

	switch l.ch {
	case '#':
		for l.peekChar() != '\n' && l.peekChar() != '\r' && l.peekChar() != 0 {
			l.readChar()
		}
		l.readChar()
		return l.NextToken()
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.RARR, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.MINUS, l.ch)
		}

	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '%':
		tok = l.newToken(token.MOD, l.ch)

	case '.':
		tok = l.newToken(token.DOT, l.ch)

	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTEQ, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LARR, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTEQ, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.GT, l.ch)
		}

	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)

	case '\r':
		if l.peekChar() == '\n' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.SEMICOLON, Literal: string(ch) + string(l.ch), Line: l.lineCount}
		} else {
			tok = l.newToken(token.SEMICOLON, l.ch)
		}
		l.lineCount++

	case '\n':
		tok = l.newToken(token.SEMICOLON, l.ch)
		l.lineCount++

	case '"':
		tok.Type = token.STRING
		x, ok := l.readString()
		if !ok {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
		tok.Literal = x

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: l.lineCount}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' { // || l.ch == '\n' || l.ch == '\r'
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	l.readChar()
	for isAlnum(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() (string, bool) {
	res := ""
	for {
		l.readChar()

		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				res += "\n"
			case 'r':
				res += "\r"
			case 't':
				res += "\t"
			case '\\':
				res += "\\"
			case '"':
				res += "\""
			case '\'':
				res += "'"
			case 'v':
				res += "\v"
			case 'a':
				res += "\a"
			case 'b':
				res += "\b"
			default:
				res += string(l.ch)
			}
			continue
		}
		if l.ch == '"' {
			return res, true
		}
		if l.ch == byte(0) {
			return "", false
		}
		res += string(l.ch)
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isAlnum(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
