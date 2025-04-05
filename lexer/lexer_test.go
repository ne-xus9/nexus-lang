package lexer

import (
	"go/token"
	"nexus/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		et token.TokenType
		el string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		token := l.NextToken()
		if token.Type != tt.et {
			t.Fatalf("tests[%d] - type is wrong, expected=%q, got=%q", i, tt.et, token.Type)
		}
		if token.Literal != tt.el {
			t.Fatalf("tests[%d] - literal is wrong, expected=%q, got=%q", i, tt.el, token.Literal)
		}
	}
}
