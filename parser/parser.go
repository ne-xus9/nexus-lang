package parser

import (
	"errors"
	"fmt"
	"nexus/ast"
	"nexus/lexer"
	"nexus/token"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > <
	SUM         // +
	PRODUCT     // *
	PREFIX      // - !
	CALL        // foo()
)

type (
	PrefixParseFn func() ast.Expression
	InfixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	lex       *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []error
	prefixFns map[token.TokenType]PrefixParseFn
	infixFns  map[token.TokenType]InfixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l, errors: []error{}}

	p.prefixFns = make(map[token.TokenType]PrefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// To populate both, curr and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) nextToken() {
	// Advance both pointers
	p.currToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{} // ast root
	prog.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RET:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken} // this looks like recursive leaves
	// This function should be called only on
	// statements which start with 'let'

	// After, requires an identifier, otherwise fails
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// Assigns an Identifier node to this statement,
	// with the Token: token.Identifier
	// and Value: p.currToken.Literal,
	// which is the identifier name
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	// After the identifier, requires a = sign,
	// if it does not exist, break and fail.
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// Keep reading until there is a semicolon
	// maybe not really good
	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// Return the statement after all.
	return stmt
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		// if peekToken is of the expected type
		// keep going and return true
		p.nextToken()
		return true
	} else {
		// Otherwise append error to parser errors
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []error {
	// Public getter for errors
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	// Append error
	msg := fmt.Sprintf("Expected token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) registerPrefix(t token.TokenType, fn PrefixParseFn) {
	p.prefixFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn InfixParseFn) {
	p.infixFns[t] = fn
}

func (p *Parser) parseExpression(prec int) ast.Expression {
	prefix := p.prefixFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse %q as integer", p.currToken)
		p.errors = append(p.errors, errors.New(msg))
	}
	lit.Value = value

	return lit
}
