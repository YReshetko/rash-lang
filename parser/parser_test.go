package parser_test

import (
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/parser"
	"github.com/YReshetko/rash-lang/tokens"
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
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
		{"add(a * b[2], import.arr[1 + func(a)], 2 * [1, 2][1])", "add((a * (b[2])), (import.(arr[(1 + func(a))])), (2 * ([1, 2][1])))"},
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

func TestSimpleReferencedExpression(t *testing.T) {
	input := `sys.time();`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	exp, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)

	left, ok := exp.Left.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "sys", left.Value)

	assert.Equal(t, ".", exp.Operator)

	right, ok := exp.Right.(*ast.CallExpression)
	require.True(t, ok)

	callIdent, ok := right.Function.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "time", callIdent.Value)
}

func TestOperationOnReferencedExpression(t *testing.T) {
	input := `let a = sys.value - 10;`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.LetStatement)
	require.True(t, ok)
	assert.Equal(t, "a", statement.Name.Value)

	infix, ok := statement.Value.(*ast.InfixExpression)
	require.True(t, ok)

	exp, ok := infix.Left.(*ast.InfixExpression)
	require.True(t, ok)

	left, ok := exp.Left.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "sys", left.Value)

	assert.Equal(t, ".", exp.Operator)

	right, ok := exp.Right.(*ast.Identifier)
	require.True(t, ok)

	assert.Equal(t, "value", right.Value)

	assert.Equal(t, "-", infix.Operator)

	intLit, ok := infix.Right.(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, int64(10), intLit.Value)
}

func TestReferencedExpression(t *testing.T) {
	input := `let s = sys.call(a, 100, fn(x, y) {return x + y;});`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.LetStatement)
	require.True(t, ok)

	exp, ok := statement.Value.(*ast.InfixExpression)
	require.True(t, ok)

	left, ok := exp.Left.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "sys", left.Value)

	assert.Equal(t, ".", exp.Operator)

	right, ok := exp.Right.(*ast.CallExpression)
	require.True(t, ok)

	callIdent, ok := right.Function.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "call", callIdent.Value)

	require.Len(t, right.Arguments, 3)

	arg1, ok := right.Arguments[0].(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "a", arg1.Value)

	arg2, ok := right.Arguments[1].(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, int64(100), arg2.Value)

	arg3, ok := right.Arguments[2].(*ast.FunctionLiteral)
	require.True(t, ok)
	assert.Len(t, arg3.Parameters, 2)
}

func TestParsingArrayLiterals(t *testing.T) {
	input := `[1, func, 2 + 3, fn(a, b){return a + b;}];`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	arr, ok := statement.Expression.(*ast.ArrayLiteral)
	require.True(t, ok)

	require.Len(t, arr.Elements, 4)

	_, ok = arr.Elements[0].(*ast.IntegerLiteral)
	assert.True(t, ok)
	_, ok = arr.Elements[1].(*ast.Identifier)
	assert.True(t, ok)
	_, ok = arr.Elements[2].(*ast.InfixExpression)
	assert.True(t, ok)
	_, ok = arr.Elements[3].(*ast.FunctionLiteral)
	assert.True(t, ok)
}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1 + 2];`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	ind, ok := statement.Expression.(*ast.IndexExpression)
	require.True(t, ok)

	left, ok := ind.Left.(*ast.Identifier)
	assert.True(t, ok)
	assert.Equal(t, "myArray", left.Value)
	infix, ok := ind.Index.(*ast.InfixExpression)
	assert.True(t, ok)
	assert.Equal(t, "(1 + 2)", infix.String())
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3};`
	expected := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	hash, ok := statement.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 3)

	for k, v := range hash.Pairs {
		strKey, ok := k.(*ast.StringLiteral)
		assert.True(t, ok)
		intVal, ok := v.(*ast.IntegerLiteral)
		assert.True(t, ok)
		exp, ok := expected[strKey.Value]
		assert.True(t, ok)
		assert.Equal(t, int64(exp), intVal.Value)
	}
}

func TestParsingHashLiteralsEmpty(t *testing.T) {
	input := `{};`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	hash, ok := statement.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 0)
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 1 + 2, "two": 2 * 3, "three": 3/1};`
	expected := map[string]string{
		"one":   "(1 + 2)",
		"two":   "(2 * 3)",
		"three": "(3 / 1)",
	}

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	hash, ok := statement.Expression.(*ast.HashLiteral)
	require.True(t, ok)
	require.Len(t, hash.Pairs, 3)

	for k, v := range hash.Pairs {
		strKey, ok := k.(*ast.StringLiteral)
		assert.True(t, ok)
		exp, ok := expected[strKey.Value]
		assert.True(t, ok)
		assert.Equal(t, exp, v.String())
	}
}

func TestAssignExpressionToLiteral(t *testing.T) {
	input := `foo = "bar";`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	infix, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)

	ident, ok := infix.Left.(*ast.Identifier)
	require.True(t, ok)
	assert.Equal(t, "foo", ident.Value)

	assert.Equal(t, "=", infix.Operator)

	value, ok := infix.Right.(*ast.StringLiteral)
	require.True(t, ok)
	assert.Equal(t, "bar", value.Value)
}

func TestAssignExpressionToIndex(t *testing.T) {
	input := `map["foo"] = "bar" + "bazz";`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	infix, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)

	_, ok = infix.Left.(*ast.IndexExpression)
	require.True(t, ok)

	assert.Equal(t, "=", infix.Operator)

	_, ok = infix.Right.(*ast.InfixExpression)
	require.True(t, ok)
}

func TestAssignExpressionToIndexFromFunc(t *testing.T) {
	input := `func()[1] = fn(){return true;}();`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	infix, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)

	_, ok = infix.Left.(*ast.IndexExpression)
	require.True(t, ok)

	assert.Equal(t, "=", infix.Operator)

	_, ok = infix.Right.(*ast.CallExpression)
	require.True(t, ok)
}

func TestAssignExpressionsExternal(t *testing.T) {
	input := `pkg.var = a + 2`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	infix, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "((pkg.var) = (a + 2))", infix.String())
}

func TestAssignExpressionsArray(t *testing.T) {
	input := `arr[10] = a + 2`

	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	infix, ok := statement.Expression.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "((arr[10]) = (a + 2))", infix.String())
}

func TestForStatement(t *testing.T) {
	input := `for (let i = 0; i < 10; i = i + 1) { sum = sum + i; }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	forExp, ok := statement.Expression.(*ast.ForExpression)
	require.True(t, ok)

	let, ok := forExp.Initial.(*ast.LetStatement)
	require.True(t, ok)
	assert.Equal(t, "let i = 0;", let.String())

	infix, ok := forExp.Condition.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(i < 10)", infix.String())

	complete, ok := forExp.Complete.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(i = (i + 1))", complete.String())

	require.Len(t, forExp.Body.Statements, 1)
	_, ok = forExp.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
}

func TestForStatementNoInitial(t *testing.T) {
	input := `for (i < 10; i = i + 1) { sum = sum + i; }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	forExp, ok := statement.Expression.(*ast.ForExpression)
	require.True(t, ok)

	assert.Nil(t, forExp.Initial)

	infix, ok := forExp.Condition.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(i < 10)", infix.String())

	complete, ok := forExp.Complete.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(i = (i + 1))", complete.String())

	require.Len(t, forExp.Body.Statements, 1)
	_, ok = forExp.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
}


func TestForStatementNoComplete(t *testing.T) {
	input := `for (i < 10;) { sum = sum + i; }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	forExp, ok := statement.Expression.(*ast.ForExpression)
	require.True(t, ok)

	assert.Nil(t, forExp.Initial)
	assert.Nil(t, forExp.Complete)

	infix, ok := forExp.Condition.(*ast.InfixExpression)
	require.True(t, ok)
	assert.Equal(t, "(i < 10)", infix.String())

	require.Len(t, forExp.Body.Statements, 1)
	_, ok = forExp.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
}


func TestForStatementNoParams(t *testing.T) {
	input := `for () { sum = sum + i; }`
	l := lexer.New(input, "non-file")
	p := parser.New(l)

	program := p.ParseProgram()
	require.Len(t, p.Errors(), 0)
	require.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	forExp, ok := statement.Expression.(*ast.ForExpression)
	require.True(t, ok)

	assert.Nil(t, forExp.Initial)
	assert.Nil(t, forExp.Complete)
	assert.Nil(t, forExp.Condition)

	require.Len(t, forExp.Body.Statements, 1)
	_, ok = forExp.Body.Statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)
}
