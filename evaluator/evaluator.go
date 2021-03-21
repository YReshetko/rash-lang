package evaluator

import (
	"fmt"
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/objects"
	"github.com/YReshetko/rash-lang/parser"
	"io/ioutil"
	"strings"
)

func Eval(node ast.Node, environment *objects.Environment) objects.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.DeclarationStatement:
		return evalDeclarationStatement(node, environment)
	case *ast.BlockStatement:
		return evalStatements(node.Statements, environment)
	case *ast.PrefixExpression:
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &objects.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolean(node.Value)
	case *ast.ReturnStatement:
		result := Eval(node.Value, environment)
		if isError(result) {
			return result
		}
		return &objects.ReturnValue{Value: result}
	case *ast.LetStatement:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}
		environment.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, environment)
	case *ast.FunctionLiteral:
		return &objects.Function{
			Parameters:  node.Parameters,
			Body:        node.Body,
			Environment: environment,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return objects.NULL
}

func evalDeclarationStatement(node *ast.DeclarationStatement, environment *objects.Environment) objects.Object {
	include, ok := node.Declaration.(*ast.IncludeDeclaration)
	if !ok {
		return newError("unknown declaration type: %s", node.Declaration.String())
	}

	src, err := ioutil.ReadFile(include.Include.Value)
	if err != nil{
		return newError("unable to load included script %s due to %v", include.Include.Value, err)
	}


	l := lexer.New(string(src), include.Include.Value)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return newError("unable to evaluate included script %s due to:\n %s", include.Include.Value, strings.Join(p.Errors(), ";\n"))
	}

	// TODO instead of filling this environment create a new one and pass it to existing for accessing it by alias
	obj := Eval(program, environment)
	if obj.Type() == objects.ERROR_OBJ {
		return obj
	}

	return objects.NULL
}

func applyFunction(function objects.Object, args []objects.Object) objects.Object {
	switch fn := function.(type) {
	case *objects.Function:
		if len(args) != len(fn.Parameters) {
			return newError("number of function parameters mismatch: expected=%d, got=%d", len(fn.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnvironment(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *objects.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", function.Type())
	}

}

func unwrapReturnValue(evaluated objects.Object) objects.Object {
	if retVal, ok := evaluated.(*objects.ReturnValue); ok {
		return retVal.Value
	}
	return evaluated
}

func extendFunctionEnvironment(fn *objects.Function, args []objects.Object) *objects.Environment {
	newEnv := objects.NewEnclosedEnvironment(fn.Environment)

	for i, parameter := range fn.Parameters {
		newEnv.Set(parameter.Value, args[i])
	}

	return newEnv
}

func evalExpressions(arguments []ast.Expression, environment *objects.Environment) []objects.Object {
	result := []objects.Object{}
	for _, argument := range arguments {
		evaluated := Eval(argument, environment)
		if isError(evaluated) {
			return []objects.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, environment *objects.Environment) objects.Object {
	if val, ok := environment.Get(node.Value); ok {
		return val
	}
	if val, ok := builtins[node.Value]; ok {
		return val
	}

	return newError("identifier not found: %s", node.Value)
}

func evalIfExpression(node *ast.IfExpression, environment *objects.Environment) objects.Object {
	condition := Eval(node.Condition, environment)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, environment)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, environment)
	}
	return objects.NULL
}

func isTruthy(obj objects.Object) bool {
	switch obj {
	case objects.NULL, objects.FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left objects.Object, right objects.Object) objects.Object {
	switch {
	case left.Type() == objects.INTEGER_OBJ && right.Type() == objects.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == objects.STRING_OBJ && right.Type() == objects.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolean(left == right)
	case operator == "!=":
		return nativeBoolean(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := left.(*objects.Integer).Value
	rightVal := right.(*objects.Integer).Value
	switch operator {
	case "+":
		return &objects.Integer{Value: leftVal + rightVal}
	case "-":
		return &objects.Integer{Value: leftVal - rightVal}
	case "/":
		return &objects.Integer{Value: leftVal / rightVal}
	case "*":
		return &objects.Integer{Value: leftVal * rightVal}
	case ">":
		return nativeBoolean(leftVal > rightVal)
	case "<":
		return nativeBoolean(leftVal < rightVal)
	case "==":
		return nativeBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left objects.Object, right objects.Object) objects.Object {
	leftVal := left.(*objects.String).Value
	rightVal := right.(*objects.String).Value
	switch operator {
	case "+":
		return &objects.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right objects.Object) objects.Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}

}

func evalMinusPrefixExpression(value objects.Object) objects.Object {
	if value.Type() != objects.INTEGER_OBJ {
		return newError("unknown operator: -%s", value.Type())
	}
	switch v := value.(type) {
	case *objects.Integer:
		return &objects.Integer{Value: -v.Value}
	default:
		return objects.NULL
	}
}

func evalBangPrefixExpression(value objects.Object) objects.Object {
	switch value {
	case objects.TRUE:
		return objects.FALSE
	case objects.FALSE:
		return objects.TRUE
	case objects.NULL:
		return objects.TRUE
	default:
		return objects.FALSE
	}
}

func evalProgram(stmts []ast.Statement, environment *objects.Environment) objects.Object {
	var result objects.Object

	for _, stmt := range stmts {
		result = Eval(stmt, environment)
		switch res := result.(type) {
		case *objects.ReturnValue:
			return res.Value
		case *objects.Error:
			res.AddStackLine(stmt.StackLine())
			return res
		}
	}
	return result
}

func evalStatements(statements []ast.Statement, environment *objects.Environment) objects.Object {
	var result objects.Object

	for _, stmt := range statements {
		result = Eval(stmt, environment)
		if result == nil || (result.Type() != objects.RETURN_VALUE_OBJ && result.Type() != objects.ERROR_OBJ) {
			continue
		}
		if result.Type() == objects.ERROR_OBJ {
			result.(*objects.Error).AddStackLine(stmt.StackLine())
		}
		return result

	}
	return result
}

func nativeBoolean(value bool) *objects.Boolean {
	if value {
		return objects.TRUE
	}
	return objects.FALSE
}

func newError(format string, args ...interface{}) *objects.Error {
	return &objects.Error{Message: fmt.Sprintf(format, args...)}
}

func isError(obj objects.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == objects.ERROR_OBJ
}
