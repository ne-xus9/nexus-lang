package parser

import (
	"nexus/ast"
	"nexus/token"
)

// ParseStatement entrypoint to parse any
// kind of statement depending on the first
// token read.
func (p *Parser) ParseStatement() ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.RET:
		return p.ParseReturnStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

// ParseLetStatement is for parsing `let foo = 30;`-like
// statements, expecting five tokens.
func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrentToken} // this looks like recursive leaves
	// This function should be called only on
	// statements which start with 'let'

	// After, requires an identifier, otherwise fails
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Assigns an Identifier node to this statement,
	// with the Token: token.Identifier
	// and Value: p.CurrentToken.Literal,
	// which is the identifier name
	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}

	// After the identifier, requires a = sign,
	// if it does not exist, break and fail.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.ParseExpression(LOWEST)

	// Keep reading until there is a semicolon
	// maybe not really good
	for !p.CurrentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Return the statement after all.
	return stmt
}

// ParseReturnStatement is for parsing `return foo;`-like
// statements, expecting only three tokens.
func (p *Parser) ParseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.CurrentToken}
	p.nextToken()

	stmt.ReturnValue = p.ParseExpression(LOWEST)

	for !p.CurrentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// ParseExpressionStatement is for parsing `foo;`-like
// statements, expecting only two tokens.
func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.CurrentToken}

	stmt.Expression = p.ParseExpression(LOWEST)

	if p.PeekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
