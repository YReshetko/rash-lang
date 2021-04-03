package evaluator_test

import (
	"github.com/YReshetko/rash-lang/evaluator"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/loaders"
	"github.com/YReshetko/rash-lang/objects"
	"github.com/YReshetko/rash-lang/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func init() {
	evaluator.ScriptLoader = loaders.ScriptLoader
}

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

func TestDoubleEval(t *testing.T) {
	tests := []struct {
		input string
		value float64
	}{
		{"5.128", 5.128},
		{"3.128 + 0.128", 3.256},
		{"3.128 - 0.128", 3},
		{"2 - -3 + 13 / 6", 7.1666666},
		{"2 - (-3 + 13) / 3", -1.3333333},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertDoubleObject(t, obj, test.value)
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
		{`
			# test "fixtures/test.rs"; 
			test.testFn();
		`, 100},
		{`
			# test "fixtures/test.rs"; 
			let sub = fn(x, y){
				return x - y;
			};
			test.anotherFn(10, 15, sub);
		`, 0},
		{`
			# test "fixtures/test.rs"; 
			let sub = fn(x, y){
				return x + y;
			};
			test.anotherFn(2, 3, sub);
		`, 10},
		{`
			# test "fixtures/test.rs"; 
			test.anotherFn(2, 3, fn(x, y){
				return x * y;
			});
		`, 12},
		{`
			# test "fixtures/test.rs"; 
			let a = test.const - 9;
			a;
		`, 120},
		{`
			# test "fixtures/test.rs";
			let valOne = test.const - 120;
			let valTwo = test.const;
			let add = fn(x, y){
				return x + y;
			};
			test.anotherFn(valOne, valTwo, add);
		`, 276},
		{`
			# test "fixtures/test.rs";
			let valOne = test.const - test.const * 2;
			let valTwo = fn(v){return v - 120}(test.const);
			let add = fn(x, y){
				return x + y;
			};
			test.anotherFn(valOne, valTwo, add);
		`, -240},
		{`
			# deep "fixtures/deep.rs";
			let add = fn(x, y){
				return x + y;
			};
			let b = deep.func(12, add);
			b;
		`, 54},
		{`
			# deep "fixtures/deepdeep.rs";
			let add = fn(x, y){
				return x + y;
			};
			let b = deep.func(add);
			b;
		`, 58},
		{
			`
			fn(a){
				if (a == 1){
					# inv "fixtures/inv_one.rs"
				} else {
					# inv "fixtures/inv_two.rs"
				}
			}(1).func()
			`, 1001,
		},
		{
			`
			fn(a){
				if (a == 1){
					# inv "fixtures/inv_one.rs"
				} else {
					# inv "fixtures/inv_two.rs"
				}
			}(2).func()
			`, 1002,
		},
		{
			`
			# inv_one "fixtures/inv_one.rs"
			# inv_two "fixtures/inv_two.rs"
			fn(a){
				if (a == 1){
					return inv_one;
				} else {
					return inv_two;
				}
			}(1).func()
			`, 1001,
		},
		{
			`
			# inv_one "fixtures/inv_one.rs"
			# inv_two "fixtures/inv_two.rs"
			fn(a){
				if (a == 1){
					return inv_one;
				} else {
					return inv_two;
				}
			}(2).func()
			`, 1002,
		},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertIntegerObject(t, obj, test.expected)
	}
}

func TestArrayDefinition(t *testing.T) {
	tests := []struct {
		input string
		value []interface{}
	}{
		{"[1, 2 * 2, 3 == 3];", []interface{}{1, 4, true}},
		{`let func = fn(a, b) {a + b}; [1, func("hello ", "world"), 3 == 3];`, []interface{}{1, "hello world", true}},
		{`# test "fixtures/test.rs"; [1, "hello " + "world", test.const];`, []interface{}{1, "hello world", 129}},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		assertArrayObject(t, obj, test.value)
	}
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"[1, 2, 3][0];", 1},
		{"[1, 2, 3][1];", 2},
		{"[1, 2, 3][2];", 3},
		{"let i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"let arr = [1, 2, 3]; arr[2];", 3},
		{"let arr = [1, 2, 3]; let i = arr[1]; arr[i];", 3},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
		{"let a = [1, 2, 3]; [a[2], a[1], a[0]]", []interface{}{3, 2, 1}},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		switch exp := test.value.(type) {
		case int:
			assertIntegerObject(t, obj, int64(exp))
		case nil:
			assertNullObject(t, obj)
		case []interface{}:
			assertArrayObject(t, obj, exp)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
		let two = "two";
		{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
			true: 5,
			false: 6
		}
`
	expected := map[objects.HashKey]int64{
		(&objects.String{"one"}).HashKey():   1,
		(&objects.String{"two"}).HashKey():   2,
		(&objects.String{"three"}).HashKey(): 3,
		(&objects.Integer{4}).HashKey():      4,
		(&objects.Boolean{true}).HashKey():   5,
		(&objects.Boolean{false}).HashKey():  6,
	}

	obj := testEval(t, input)
	hash, ok := obj.(*objects.Hash)
	require.True(t, ok)
	for key, value := range expected {
		pair, ok := hash.Pairs[key]
		assert.True(t, ok)
		assert.Equal(t, value, pair.Value.(*objects.Integer).Value)
	}

}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "bar"; {"foo": 5}[key]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
		{`{"arr": [1, 2, 3]}["arr"]`, []interface{}{1, 2, 3}},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		switch exp := test.value.(type) {
		case int:
			assertIntegerObject(t, obj, int64(exp))
		case nil:
			assertNullObject(t, obj)
		case []interface{}:
			assertArrayObject(t, obj, exp)
		}
	}
}

func TestAssignExpression(t *testing.T) {
	tests := []struct {
		input string
		value interface{}
	}{
		{"let a = 5; a = 10; a;", 10},
		{`let a = ["hello", "world"]; a[1] = "rash"; a;`, []interface{}{"hello", "rash"}},
		{`let m = {true: false, false: true}; m[true] = true; m[true];`, true},
		{`let func = fn(){return {"one":"hello","two": "world"};}; let b = {}; b = func(); b["two"];`, "world"},
		{`let m = {"one":"hello","two": "world"}; let func = fn(){return m;}; func()["two"] = "rash"; m["two"];`, "rash"},
	}
	for _, test := range tests {
		obj := testEval(t, test.input)
		switch exp := test.value.(type) {
		case int:
			assertIntegerObject(t, obj, int64(exp))
		case string:
			assertStringObject(t, obj, exp)
		case []interface{}:
			assertArrayObject(t, obj, exp)
		case bool:
			assertBooleanObject(t, obj, exp)
		}

	}
}

func TestForExpression(t *testing.T) {
	input := `
		let sum = 0;
		for (let i = 0; i < 10; i = i + 1)
		{
			sum = sum + i
		}
		sum;
`
	obj := testEval(t, input)
	assertIntegerObject(t, obj, 45)
}

func TestForExpressionCondition(t *testing.T) {
	input := `
		let sum = 0;
		for ( sum < 10; sum = sum + 1){}
		sum;
`
	obj := testEval(t, input)
	assertIntegerObject(t, obj, 10)
}

func TestForExpressionConditionOnly(t *testing.T) {
	input := `
		let sum = 0;
		for ( sum < 15;){
			sum = sum + 1;
		}
		sum;
`
	obj := testEval(t, input)
	assertIntegerObject(t, obj, 15)
}

func TestForExpressionEmpty(t *testing.T) {
	input := `
		let sum = 0;
		fn() {
			for (){
				sum = sum + 1;
				if (sum == 5) {
					return;
				}
			}
		}()
		sum;
`
	obj := testEval(t, input)
	assertIntegerObject(t, obj, 5)
}

func TestForAssignment(t *testing.T) {
	input := `
		let sum = 0;
		let a = 1;
		a = for (){
			sum = sum + 1;
			if (sum == 5) {
				return;
			}
		}
		sum;
`
	obj := testEval(t, input)
	assertIntegerObject(t, obj, 5)
}

func assertArrayObject(t *testing.T, obj objects.Object, value []interface{}) {
	arr, ok := obj.(*objects.Array)
	require.True(t, ok)
	require.Equal(t, len(value), len(arr.Elements))
	for i, v := range value {
		switch exp := v.(type) {
		case int:
			assertIntegerObject(t, arr.Elements[i], int64(exp))
		case string:
			assertStringObject(t, arr.Elements[i], exp)
		case bool:
			assertBooleanObject(t, arr.Elements[i], exp)
		default:
			assert.Fail(t, "unsupported type %T", v)
		}
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

func assertDoubleObject(t *testing.T, obj objects.Object, expectedValue float64) {
	require.NotNil(t, obj)
	assert.Equal(t, objects.DOUBLE_OBJ, obj.Type())
	intObj, ok := obj.(*objects.Double)
	require.True(t, ok)
	assert.True(t, math.Abs(expectedValue-intObj.Value) < 0.000001)
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
