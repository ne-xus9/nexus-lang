package lexer

import (
	"nexus/token"
)

type Lexer struct {
	input   string // Code source
	pos     uint   // Buffer position
	readPos uint   // Right limiter
	ch      byte   // Actual character
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		if l.peekNext() == '=' {
			ch := l.ch
			l.readChar()
			t = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			t = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		t = newToken(token.SEMICOLON, l.ch)
	case '(':
		t = newToken(token.LPAREN, l.ch)
	case ')':
		t = newToken(token.RPAREN, l.ch)
	case ',':
		t = newToken(token.COMMA, l.ch)
	case '{':
		t = newToken(token.LBRACE, l.ch)
	case '}':
		t = newToken(token.RBRACE, l.ch)
	case '+':
		t = newToken(token.PLUS, l.ch)
	case '-':
		t = newToken(token.SUBS, l.ch)
	case '*':
		t = newToken(token.MULT, l.ch)
	case '/':
		t = newToken(token.DIV, l.ch)
	case '<':
		t = newToken(token.LT, l.ch)
	case '>':
		t = newToken(token.GT, l.ch)
	case '!':
		if l.peekNext() == '=' {
			ch := l.ch
			l.readChar()
			t = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			t = newToken(token.NOT, l.ch)
		}
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.ch) {
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)
			return t
		} else if isDigit(l.ch) {
			t.Type = token.INT
			t.Literal = l.readNumber()
			return t
		} else {
			t = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return t
}
