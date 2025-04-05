package lexer

import "nexus/token"

type Lexer struct {
	input   string
	pos     uint
	readPos uint
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= uint(len(l.input)) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func newToken(tt token.TokenType, ch byte) token.Token {
	return token.Token{Type: tt, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	switch l.ch {
	case '=':
		t = newToken(token.ASSIGN, l.ch)
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
	case 0:
		t.Literal = ""
		t.Type = token.EOF
	}
	l.readChar()
	return t
}
