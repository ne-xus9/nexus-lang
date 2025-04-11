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
		t.Fatalf("Could not build a program with input '%s'", input)
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Wrong number of statements. got=%d", len(program.Statements))
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

func TestIdentifiers(t *testing.T) {
	parser := New(lexer.New("foobar;"))

	prog := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Statements has wrong length. got=%d", len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is not an ast.ExpressionStatement. got=%T", prog.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt is not an ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident.Value is not foobar. got=%s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
		5;
		foobar;
		x;
	`
	parser := New(lexer.New(input))
	prog := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(prog.Statements) != 3 {
		t.Fatalf("prog.Statements has wrong length. got=%d", len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is not an ast.ExpressionStatement. got=%T", prog.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("Literal.Value is not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Literal.TokenLiteral is not %s. got=%s", "5", literal.TokenLiteral())
	}
}
