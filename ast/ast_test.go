package ast

import (
	"nexus/token"
	"testing"
)

func TestAsString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Value: "foo",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "bar"},
					Value: "bar",
				},
			},
		},
	}
	if prog.AsString() != "let foo = bar;" {
		t.Errorf("prog.AsString() wrong. got=%q", prog.AsString())
	}
}
