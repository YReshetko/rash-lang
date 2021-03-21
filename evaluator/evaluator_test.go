package evaluator_test

import (
	"github.com/YReshetko/monkey-language/evaluator"
	"github.com/YReshetko/monkey-language/lexer"
	"github.com/YReshetko/monkey-language/objects"
	"github.com/YReshetko/monkey-language/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntegerEval(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"5", 5},
		{"132312", 132312},
		{"-5", -5},
		{"-132312", -132312},
		{"2 + 10 - 12", 0},
		{"2 * (10 - 12)", -4},
		{"5 - 3 * 2 + 4", 3},
		{"2 - -3 + 12 / 6", 7},
		{"2 - (-3 + 12) / 3", -1},
		{"50 - 100 + 50", 0},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.value)
	}
}

func TestStringEval(t *testing.T) {
	tests := []struct {
		input string
		value string
	}{
		{`"hello world";`, "hello world"},
		{`let a = "hello world"; a;`, "hello world"},
		{`let a = "hello" + " " + "world"; a;`, "hello world"},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertStringObject(t, obj, test.value)
	}
}

func TestBooleanEval(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"true", true},
		{"false", false},
		{"2 != 2", false},
		{"2 != 3", true},
		{"100 < 49", false},
		{"100 > 49", true},
		{"100 == 49", false},
		//{"100 > 40 > 30", true},
		{"true == true", true},
		{"true != true", false},
		{"false != false", false},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false == true", false},
		{"false != true", true},
		{"(1>2) == true", false},
		{"(1>2) != true", true},
		{`"hello" == "world"`, false},
		{`"hello" != "world"`, true},
		{`"hello" == "hello"`, true},
		{`"hello" != "hello"`, false},
		{`let a = "hello"; a != "hello";`, false},
		{`let a = "hello"; a == "hello";`, true},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertBooleanObject(t, obj, test.value)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		value bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!false", false},
		{"!!true", true},
		{"!!!true", false},
		{"!!!false", true},
		{"!5", false},
		{"!!5", true},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertBooleanObject(t, obj, test.value)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) {10}", 10},
		{"if (false) {10}", nil},
		{"if (2 > 10) {10}", nil},
		{"if (2 < 10) {10}", 10},
		{"if (2 != 10) {10}", 10},
		{"if (2 == 10) {10} else {20}", 20},
		{
			` if (true) {
						if (true) {
							return 10;
						}
						return 1;
					}`, 10,
		},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		integer, ok := test.expected.(int)
		if ok {
			assertIntegerObject(t, obj, int64(integer))
		} else {
			assertNullObject(t, obj)
		}

	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"9; return 10; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.value)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"false - 5; 5;",
			"type mismatch: BOOLEAN - INTEGER",
		},
		{
			"-true;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 2) {true + false;}",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 2) { "Hello" - "world";}`,
			"unknown operator: STRING - STRING",
		},
		{
			` if (true) {
						if (true) {
							return true + false;
						}
						return 1;
					}`, "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
	}

	for _, test := range tests {
		obj := testEval(t, test.input)
		assertError(t, obj, test.error)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a * 5; b;", 25},
		{"let a = 5; let b = 2 * a; b;", 10},
		{"let a = 5; let b = 2 * a; let c = a + b * 2; c;", 25},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.value)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x){x + 2;};"

	obj := testEval(t, input)
	require.NotNil(t, obj)
	require.Equal(t, objects.FUNCTION_OBJ, obj.Type())

	fnObj, ok := obj.(*objects.Function)
	require.True(t, ok)

	require.Len(t, fnObj.Parameters, 1)
	assert.Equal(t, "x", fnObj.Parameters[0].String())

	assert.Equal(t, "(x + 2)", fnObj.Body.String())
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		value int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(a) { return a * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 3);", 8},
		{"let add = fn(x, y) { return x + y; }; add(5 + 2, add(5, 5));", 17},
		{"fn(x, y) { return x + y; }(5, 10);", 15},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.value)
	}
}

func TestIncludeDeclarations(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`# test "fixtures/test.rs"; testFn();`, 100},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.expected)
	}
}

func assertError(t *testing.T, obj objects.Object, e string) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.ERROR_OBJ, obj.Type())
	errObj, ok := obj.(*objects.Error)
	require.True(t, ok)
	assert.Equal(t, e, errObj.Message)
}

func assertNullObject(t *testing.T, obj objects.Object) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.NULL, obj)
}

func assertIntegerObject(t *testing.T, obj objects.Object, expectedValue int64) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.INTEGER_OBJ, obj.Type())
	intObj, ok := obj.(*objects.Integer)
	require.True(t, ok)
	assert.Equal(t, expectedValue, intObj.Value)
}

func assertBooleanObject(t *testing.T, obj objects.Object, expectedValue bool) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.BOOLEAN_OBJ, obj.Type())
	boolObj, ok := obj.(*objects.Boolean)
	require.True(t, ok)
	assert.Equal(t, expectedValue, boolObj.Value)
}

func assertStringObject(t *testing.T, obj objects.Object, expectedValue string) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.STRING_OBJ, obj.Type())
	strObj, ok := obj.(*objects.String)
	require.True(t, ok)
	assert.Equal(t, expectedValue, strObj.Value)
}

func testEval(t *testing.T, input string) objects.Object {
	l := lexer.New(input, "non-file")
	require.NotNil(t, l)

	p := parser.New(l)
	require.NotNil(t, p)

	node := p.ParseProgram()
	require.Empty(t, p.Errors())

	return evaluator.Eval(node, objects.NewEnvironment())
}
