package parser

import "nexus/ast"

func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}
