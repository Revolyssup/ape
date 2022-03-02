package parser

import (
	"fmt"
	"testing"

	"github.com/Revolyssup/ape/ast"
	"github.com/Revolyssup/ape/lexer"
)

//This function runs on the premise that there were no errors encountered in checkParserErrors
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}
	letstmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letstmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letstmt.Name.Value)

		return false
	}

	if letstmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letstmt.Name)
		return false
	}
	return true
}
func TestLetStatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		fmt.Println(i+1, stmt.TokenLiteral(), "statement")
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return ass();
	return 1;
	return 2;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	for _, stmt := range program.Statements {
		rstmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement not ast.Returnstatement, got=%s", stmt)
			continue
		}
		if rstmt.TokenLiteral() != "return" {
			t.Errorf("statement not ast.Returnstatement, got=%q", rstmt.TokenLiteral())
		}
	}

}

func TestExpression_IDENTIFIER_Statement(t *testing.T) {
	input := `ident;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.statements[0] is not ast.ExpressionStatement. got = %T", program.Statements[0])
	}

	//testing for identifier
	identstmt, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected ast.Identifier, got = %T", identstmt)
	}

	if identstmt.Value != "ident" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", identstmt.Value)

	}
	if identstmt.TokenLiteral() != "ident" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			identstmt.TokenLiteral())
	}

}

func TestExpression_INTEGER_LITERAL_Statement(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.statements[0] is not ast.ExpressionStatement. got = %T", program.Statements[0])
	}

	//testing for integer literal
	intstmt, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected ast.Integer, got = %T", intstmt)
	}

	if intstmt.Value != 5 {
		t.Fatalf("Expected 5, got = %q", intstmt.Value)
	}
}

func TestExpression_PREFIX(t *testing.T) {
	test := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-6;", "-", 6},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range test {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

	}

}
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("Program has %d errors", len(errors))

	for _, msg := range errors {
		t.Errorf("[Monke angry: ] %q ", msg)
	}
	t.FailNow()
}

//helper
func testIntegerLiteral(t *testing.T, il ast.Expression, value interface{}) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}
	return true

}
func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.LeftExpression, tt.leftValue) {
			return
		}
		if !testLiteralExpression(t, exp.RightExpression, tt.rightValue) {
			return
		}

	}
}
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}
	return true
}
func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, bool(v))
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"2*3+5 ",
			"((2 * 3) + 5)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x } else {y}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}
	fmt.Println(exp.String())
	if len(exp.MainStmt.Stmts) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.MainStmt.Stmts))
	}
	consequence, ok := exp.AltStmt.Stmts[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.MainStmt.Stmts[0])
	}

	if !testIdentifier(t, consequence.Expression, "y") {
		return
	}
	if exp.AltStmt == nil {
		t.Errorf("exp.Alternative.Statements was  nil. got=%+v", exp.AltStmt)
	}
}
func TestForExpression(t *testing.T) {
	input := `for (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T",
			stmt.Expression)
	}
	fmt.Println(exp.String())
	if len(exp.Stmt.Stmts) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Stmt.Stmts))
	}

}
func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}
	if len(function.Params) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Params))
	}
	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")
	if len(function.Body.Stmts) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Stmts))
	}
	// bodyStmt, ok := function.Body.Stmts[0].(*ast.ExpressionStatement)
	// if !ok {
	// 	t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
	// 		function.Body.Stmts[0])
	// }
	// testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}
func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestArray(t *testing.T) {
	input := `[5,1,12]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if literal.String() != `[5,1,12]` {
		t.Errorf("literal.String() not %q. got=%q", `[5,1,12]`, literal.String())
	}
}

// func TestObject(t *testing.T) {
// 	input := `{
// 		"name:"Ashish",
// 		"roll":2
// 	}`
// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)
// 	stmt := program.Statements[0].(*ast.ExpressionStatement)
// 	literal, ok := stmt.Expression.(*ast.ObjectLiteral)
// 	if !ok {
// 		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
// 	}
// 	if literal.String() != `[5,1,12,]` {
// 		t.Errorf("literal.String() not %q. got=%q", `[5,1,12,]`, literal.String())
// 	}
// }
func TestArrayEle(t *testing.T) {
	input := `a[0]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.ArrObjElement)
	if !ok {
		t.Fatalf("exp not *ast.ArrEle. got=%T", stmt.Expression)
	}
	if literal.String() != `a[0]` {
		t.Errorf("literal.String() not %q. got=%q", `a[0]`, literal.String())
	}
}
