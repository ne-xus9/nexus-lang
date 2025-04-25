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

	p.infixFns = make(map[token.TokenType]InfixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.SUBS, p.parseInfixExpression)
	p.registerInfix(token.MULT, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

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

var precedences = map[token.TokenType]int{
	token.EQ:   EQUALS,
	token.NEQ:  EQUALS,
	token.LT:   LESSGREATER,
	token.GT:   LESSGREATER,
	token.PLUS: SUM,
	token.SUBS: SUM,
	token.DIV:  PRODUCT,
	token.MULT: PRODUCT,
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.CurrentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.PeekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.CurrentToken,
		Operator: p.CurrentToken.Literal,
		Left:     left,
	}
	prec := p.currPrecedence()
	p.nextToken()
	expression.Right = p.ParseExpression(prec)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.CurrentToken, Value: p.CurrentTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.ParseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.CurrentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expression.Condition = p.ParseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.PeekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.CurrentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.CurrentTokenIs(token.RBRACE) && !p.CurrentTokenIs(token.EOF) {
		stmt := p.ParseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.CurrentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.PeekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.PeekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}
