package parser

import (
	"nexus/ast"
	"nexus/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foo = 8183910;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.FailNow()
	}
	if len(program.Statements) != 3 {
		t.FailNow()
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetSatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error: %q", e.Error())
	}
	t.FailNow()
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return variable;
		return 0;
	`
	lex := lexer.New(input)
	parser := New(lex)
	prog := parser.ParseProgram()

	checkParserErrors(t, parser)
	if len(prog.Statements) != 3 {
		t.Fatalf("prog.Statements does not contain 3 statements. got=%d", len(prog.Statements))
	}

	for _, val := range prog.Statements {
		stmt, ok := val.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatment. got=%T", stmt)
		}

		if stmt.TokenLiteral() != "return" {
			t.Errorf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
		}
	}
}

func TestConstStatement(t *testing.T) {
	input := `
		const MINUTES = 5;
		const SECRET_NUM = 20;
		const __PRIVATE_CONST = 0;
	`
	lex := lexer.New(input)
	parser := New(lex)
	prog := parser.ParseProgram()

	checkParserErrors(t, parser)
	if len(prog.Statements) != 3 {
		t.Fatalf("prog.Statements does not contain 3 statements. got=%d", len(prog.Statements))
	}

	for _, val := range prog.Statements {
		stmt, ok := val.(*ast.ConstStatement)
		if !ok {
			t.Errorf("stmt not *ast.ConstStatement. got=%T", stmt)
		}

		if stmt.TokenLiteral() != "const" {
			t.Errorf("stmt.TokenLiteral not 'const'. got=%q", stmt.TokenLiteral())
		}
	}
}
