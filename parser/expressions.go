package parser

import "nexus/ast"

// ParseExpression is for parsing `foo;`-like
// expressions, helper for the actual
// ParseExpressionStatement function.
func (p *Parser) ParseExpression(prec int) ast.Expression {
	prefix := p.prefixFns[p.CurrentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}
