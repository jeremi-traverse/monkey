package parser

import (
	"fmt"
	"testing"

	"github.com/jeremi-traverse/monkey/ast"
	"github.com/jeremi-traverse/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("progam.Statements does not contain 3 statments. got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{

		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 10;
		return y;
		return 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("progam.Statements does not contain 3 statments. got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected Return Statement got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Expected token literal 'return' got %s instead",
				returnStmt.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, expectedName string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("Statement isn't a LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != expectedName {
		t.Errorf("letStmt.Name not %q. got=%q", expectedName, letStmt.Name)
		return false
	}

	if letStmt.Name.TokenLiteral() != expectedName {
		t.Errorf("letStmt.Name.TokenLiteral() not %q. got=%q",
			expectedName,
			letStmt.Name.TokenLiteral())
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("statement is not an Identifier. got=%T",
			ident)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident.Value is not %s . got=%s",
			"foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() is not %s . got=%s",
			"foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("statement is not an *IntegerLiteral. got=%T",
			stmt.Expression)
	}

	if ident.Value != 5 {
		t.Fatalf("ident.Value is not %d . got=%d",
			5, ident.Value)
	}

	if ident.TokenLiteral() != "5" {
		t.Fatalf("ident.TokenLiteral() is not %s . got=%s",
			"5", ident.TokenLiteral())
	}

}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				stmt.Expression)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T",
				stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s. got=%s", exp.Operator, tt.operator)
		}

		if testInteger(t, tt.integerValue, exp.Right) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input        string
		leftOperand  int64
		operator     string
		rightOperand int64
	}{
		{"5 - 5", 5, "-", 5},
		{"5 + 5", 5, "+", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"5 < 5", 5, "<", 5},
		{"5 > 5", 5, ">", 5},
	}

	for _, tc := range infixTests {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("len(program.Statements) is not 1. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] isn't *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		infixExp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt isn't *ast.InfixExpression. got=%T",
				stmt.Expression)

		}

		if !testInteger(t, tc.leftOperand, infixExp.Left) {
			return
		}

		if infixExp.Operator != tc.operator {
			t.Fatalf("operator is not %s. got=%s", tc.operator, infixExp.Operator)
		}

		if !testInteger(t, tc.rightOperand, infixExp.Right) {
			return
		}
	}
}

func testInteger(t *testing.T, expected int64, got ast.Expression) bool {
	integer, ok := got.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("got is not an *ast.IntegerLiteral, is=%T", got)
		return false
	}

	if integer.Value != expected {
		t.Fatalf("integer.Value isn't %d. got=%d", expected, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Fatalf("integer.TokenLiteral() isn't %d. got=%s",
			expected, integer.TokenLiteral())
		return false
	}

	return true
}
