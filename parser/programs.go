package parser

import (
	"nexus/ast"
	"nexus/token"
)

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{} // ast root
	prog.Statements = []ast.Statement{}

	for p.CurrentToken.Type != token.EOF {
		if stmt := p.ParseStatement(); stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}

	return prog
}
