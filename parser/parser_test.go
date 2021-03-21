package parser_test

import (
	"github.com/YReshetko/monkey-language/ast"
	"github.com/YReshetko/monkey-language/lexer"
	"github.com/YReshetko/monkey-language/parser"
	"github.com/YReshetko/monkey-language/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let value = 8345;
`
	l := lexer.New(input, "non-file")

	p := parser.New(l)
	require.NotNil(t, p)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)

	tests := []struct {
		expectedIdentifier string
		expectedValue      string
	}{
		{"x", "5"},
		{"y", "10"},
		{"value", "8345"},
	}

	for i, test := range tests {
		statement := program.Statements[i]
		assert.Equal(t, "let", statement.TokenLiteral())

		letStatement, ok := statement.(*ast.LetStatement)
		assert.True(t, ok)

		assert.Equal(t, test.expectedIdentifier, letStatement.Name.Value)
		assert.Equal(t, test.expectedIdentifier, letStatement.Name.TokenLiteral())
		assert.Equal(t, test.expectedValue, letStatement.Value.String())
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return add(a,   b);
`
	values := []string{"5", "10", "add(a, b)"}
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)

	for i, statement := range program.Statements {
		assert.NotNil(t, statement)
		returnStatement, ok := statement.(*ast.ReturnStatement)
		assert.True(t, ok)
		assert.Equal(t, "return", returnStatement.TokenLiteral())
		assert.Equal(t, values[i], returnStatement.Value.String())
	}

}

func TestIdentifier(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	identifier, ok := statement.Expression.(*ast.Identifier)
	require.True(t, ok)

	assert.Equal(t, "foobar", identifier.TokenLiteral())
	assert.Equal(t, "foobar", identifier.Value)
}

func TestInteger(t *testing.T) {
	input := `5;`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	identifier, ok := statement.Expression.(*ast.IntegerLiteral)
	require.True(t, ok)

	assert.Equal(t, tokens.TokenType(tokens.INT), identifier.Token.Type)
	assert.Equal(t, "5", identifier.TokenLiteral())
	assert.Equal(t, int64(5), identifier.Value)
}


func TestStringLiteral(t *testing.T) {
	input := `"hello \"world\"";`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	strLit, ok := statement.Expression.(*ast.StringLiteral)
	require.True(t, ok)

	assert.Equal(t, tokens.TokenType(tokens.STRING), strLit.Token.Type)
	assert.Equal(t, `hello "world"`, strLit.Value)
}


func TestPrefixExpressionsInteger(t *testing.T) {
	tests := []struct {
		input            string
		expectedOperator string
		expectedValue    int64
	}{
		{"!5;", "!", 5},
		{"-10", "-", 10},
	}
	for _, test := range tests {

		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		require.True(t, ok)

		assert.Equal(t, test.expectedOperator, exp.Operator)

		intLiteral, ok := exp.Right.(*ast.IntegerLiteral)
		require.True(t, ok)
		assert.Equal(t, test.expectedValue, intLiteral.Value)
		assert.Equal(t, strconv.Itoa(int(test.expectedValue)), intLiteral.TokenLiteral())

	}
}

func TestPrefixExpressionsIdentifier(t *testing.T) {
	tests := []struct {
		input              string
		expectedOperator   string
		expectedIdentifier string
	}{
		{"!x", "!", "x"},
		{"-myVar", "-", "myVar"},
	}
	for _, test := range tests {

		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		exp, ok := statement.Expression.(*ast.PrefixExpression)
		require.True(t, ok)

		assert.Equal(t, test.expectedOperator, exp.Operator)

		intLiteral, ok := exp.Right.(*ast.Identifier)
		require.True(t, ok)
		assert.Equal(t, test.expectedIdentifier, intLiteral.Value)
		assert.Equal(t, test.expectedIdentifier, intLiteral.TokenLiteral())
	}
}

func TestInfixExpressionsInteger(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 - 5", 5, "-", 5},
		{"5 + 5", 5, "+", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
	}
	for _, test := range tests {

		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		exp, ok := statement.Expression.(*ast.InfixExpression)
		require.True(t, ok)

		assert.Equal(t, test.operator, exp.Operator)

		intLiteral, ok := exp.Left.(*ast.IntegerLiteral)
		require.True(t, ok)
		assert.Equal(t, test.rightValue, intLiteral.Value)
		assert.Equal(t, strconv.Itoa(int(test.rightValue)), intLiteral.TokenLiteral())

		intLiteral, ok = exp.Right.(*ast.IntegerLiteral)
		require.True(t, ok)
		assert.Equal(t, test.rightValue, intLiteral.Value)
		assert.Equal(t, strconv.Itoa(int(test.rightValue)), intLiteral.TokenLiteral())

	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"-a", "(-a)"},
		{"!5", "(!5)"},
		{"a + b", "(a + b)"},
		{"a + -5", "(a + (-5))"},
		{"a + -5 - -3", "((a + (-5)) - (-3))"},
		{"5 * a - b + -a * -10 == 35 * -myVar + 23 / 34 + 2 * -myVar", "((((5 * a) - b) + ((-a) * (-10))) == (((35 * (-myVar)) + (23 / 34)) + (2 * (-myVar))))"},
		{"(a + -5) * -3", "((a + (-5)) * (-3))"},
		{"a + -5 * -3", "(a + ((-5) * (-3)))"},
		{"a * 2 / (3 + -2) * 5", "(((a * 2) / (3 + (-2))) * 5)"},
		{"a * 2 / (3 + add(2 - b, sub(3), b + sq(14 - 3)) - -2) * 5", "(((a * 2) / ((3 + add((2 - b), sub(3), (b + sq((14 - 3))))) - (-2))) * 5)"},
	}

	for _, test := range tests {

		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		assert.Equal(t, test.output, statement.String())
	}
}

func TestBooleanExpressions(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"true", "true"},
		{"false", "false"},
		{"a > b == true", "((a > b) == true)"},
		{"!false != !a > b", "((!false) != ((!a) > b))"},
	}

	for _, test := range tests {

		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		assert.Equal(t, test.output, statement.String())
	}
}

func TestIfStatement(t *testing.T) {
	input := `if (x > y) { x }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	ifExp, ok := statement.Expression.(*ast.IfExpression)
	require.True(t, ok)

	infix, ok := ifExp.Condition.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(x > y)", infix.String())

	require.Len(t, ifExp.Consequence.Statements, 1)
	consStatement, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	consIdent, ok := consStatement.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "x", consIdent.Value)

	assert.Nil(t, ifExp.Alternative)

}

func TestIfElseStatement(t *testing.T) {
	input := `if (x > y) { x } else { y }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	ifExp, ok := statement.Expression.(*ast.IfExpression)
	require.True(t, ok)

	infix, ok := ifExp.Condition.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(x > y)", infix.String())

	require.Len(t, ifExp.Consequence.Statements, 1)
	consStatement, ok := ifExp.Consequence.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	consIdent, ok := consStatement.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "x", consIdent.Value)

	require.Len(t, ifExp.Alternative.Statements, 1)
	alterStatement, ok := ifExp.Alternative.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	alterIdent, ok := alterStatement.Expression.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "y", alterIdent.Value)

}

func TestFunctionLiteral(t *testing.T) {
	input := `fn (x, y) { x + y;}`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	fnLit, ok := statement.Expression.(*ast.FunctionLiteral)
	require.True(t, ok)

	assert.Len(t, fnLit.Parameters, 2)
	assert.Equal(t, "x", fnLit.Parameters[0].Value)
	assert.Equal(t, "y", fnLit.Parameters[1].Value)

	assert.Len(t, fnLit.Body.Statements, 1)

	expStmt, ok := fnLit.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)
	infExp, ok := expStmt.Expression.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(x + y)", infExp.String())
}

func TestFunctionLiteral_Parameters(t *testing.T) {
	tests := []struct {
		input  string
		params []string
	}{
		{"fn(){}", []string{}},
		{"fn(a){}", []string{"a"}},
		{"fn(a,   b){}", []string{"a", "b"}},
		{"fn(a, myVar,  b){}", []string{"a", "myVar", "b"}},
	}
	for _, test := range tests {
		l := lexer.New(test.input, "non-file")
		p := parser.New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		fnLit, ok := statement.Expression.(*ast.FunctionLiteral)
		require.True(t, ok)

		assert.Len(t, fnLit.Parameters, len(test.params))
		for i, param := range test.params {
			assert.Equal(t, param, fnLit.Parameters[i].Value)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add (a + 5, 2 * 3, 3)`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	call, ok := statement.Expression.(*ast.CallExpression)
	require.True(t, ok)

	ident, ok := call.Function.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "add", ident.TokenLiteral())

	require.Len(t, call.Arguments, 3)
	assert.Equal(t, "(a + 5)", call.Arguments[0].String())
	assert.Equal(t, "(2 * 3)", call.Arguments[1].String())
	assert.Equal(t, "3", call.Arguments[2].String())
}

func TestDeclarationStatement(t *testing.T) {
	input := `
# sys "sys.rs";
# time "path/to/script/time.rs"
`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 2)

	statement, ok := program.Statements[0].(*ast.DeclarationStatement)
	require.True(t, ok)

	include, ok := statement.Declaration.(*ast.IncludeDeclaration)
	require.True(t, ok)

	assert.Equal(t, tokens.TokenType(tokens.HASH), include.Token.Type)
	assert.Equal(t, "sys", include.Alias.Value)
	assert.Equal(t, "sys.rs", include.Include.Value)


	statement, ok = program.Statements[1].(*ast.DeclarationStatement)
	require.True(t, ok)

	include, ok = statement.Declaration.(*ast.IncludeDeclaration)
	require.True(t, ok)

	assert.Equal(t, tokens.TokenType(tokens.HASH), include.Token.Type)
	assert.Equal(t, "time", include.Alias.Value)
	assert.Equal(t, "path/to/script/time.rs", include.Include.Value)
}