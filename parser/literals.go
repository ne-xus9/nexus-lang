package parser

import (
	"errors"
	"fmt"
	"nexus/ast"
	"strconv"
)

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.CurrentToken}

	value, err := strconv.ParseInt(p.CurrentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse %q as integer", p.CurrentToken)
		p.errors = append(p.errors, errors.New(msg))
	}
	lit.Value = value

	return lit
}
