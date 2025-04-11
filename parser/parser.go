package parser

import (
	"errors"
	"fmt"
	"nexus/ast"
	"nexus/lexer"
	"nexus/token"
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
	lex          *lexer.Lexer
	CurrentToken token.Token
	PeekToken    token.Token
	errors       []error
	prefixFns    map[token.TokenType]PrefixParseFn
	infixFns     map[token.TokenType]InfixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l, errors: []error{}}

	p.prefixFns = make(map[token.TokenType]PrefixParseFn)
	p.registerPrefix(token.IDENT, p.ParseIdentifier)
	p.registerPrefix(token.INT, p.ParseIntegerLiteral)
	p.registerPrefix(token.NOT, p.ParsePrefixExpression)
	p.registerPrefix(token.SUBS, p.ParsePrefixExpression)

	// To populate both, curr and peek tokens
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	// Advance both pointers
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lex.NextToken()
}

func (p *Parser) CurrentTokenIs(t token.TokenType) bool {
	return p.CurrentToken.Type == t
}

func (p *Parser) PeekTokenIs(t token.TokenType) bool {
	return p.PeekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.PeekTokenIs(t) {
		// if PeekToken is of the expected type
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
	msg := fmt.Sprintf("Expected token to be %s, got %s instead", t, p.PeekToken.Type)
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) registerPrefix(t token.TokenType, fn PrefixParseFn) {
	p.prefixFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn InfixParseFn) {
	p.infixFns[t] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, errors.New(msg))
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.CurrentToken,
		Operator: p.CurrentToken.Literal,
	}
	p.nextToken()

	exp.Right = p.ParseExpression(PREFIX)
	return exp
}
