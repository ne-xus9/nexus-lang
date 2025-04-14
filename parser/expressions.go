package parser

import (
	"nexus/ast"
	"nexus/token"
)

// ParseExpression is for parsing `foo;`-like
// expressions, helper for the actual
// ParseExpressionStatement function.
func (p *Parser) ParseExpression(prec int) ast.Expression {
	prefix := p.prefixFns[p.CurrentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.CurrentToken.Type)
		return nil
	}
	leftExp := prefix()

	if !p.PeekTokenIs(token.SEMICOLON) && prec < p.peekPrecedence() {
		infix := p.infixFns[p.PeekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
