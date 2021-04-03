package evaluator

import (
	"errors"
	"fmt"
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/objects"
)

type scriptLoader func(path string) (*objects.Environment, error)

var ScriptLoader scriptLoader = func(path string) (*objects.Environment, error) {
	return nil, errors.New("script loader is not defined")
}

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
		return evalInfixExpression(node, environment)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.ForExpression:
		return evalForExpression(node, environment)
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}
	case *ast.DoubleLiteral:
		return &objects.Double{Value: node.Value}
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
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, environment)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &objects.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, environment)
	case *ast.IndexExpression:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, environment)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}

	return objects.NULL
}

func evalHashLiteral(node *ast.HashLiteral, environment *objects.Environment) objects.Object {
	hash := &objects.Hash{
		Pairs: map[objects.HashKey]objects.HashPair{},
	}

	for keyExp, valueExp := range node.Pairs {
		key := Eval(keyExp, environment)
		if isError(key) {
			return key
		}
		hashable, ok := key.(objects.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueExp, environment)
		if isError(value) {
			return value
		}

		hash.Pairs[hashable.HashKey()] = objects.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return hash
}

func evalIndexExpression(left objects.Object, index objects.Object) objects.Object {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == objects.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported for: %s", left.Type())
	}
}

func evalHashIndexExpression(left objects.Object, index objects.Object) objects.Object {
	hash := left.(*objects.Hash)
	ind, ok := index.(objects.Hashable)
	if !ok {
		return newError("unusable as a hash key: %s", index.Type())
	}
	pair, ok := hash.Pairs[ind.HashKey()]
	if !ok {
		return objects.NULL
	}
	return pair.Value
}

func evalArrayIndexExpression(left objects.Object, index objects.Object) objects.Object {
	arr := left.(*objects.Array)
	ind := index.(*objects.Integer).Value
	max := int64(len(arr.Elements) - 1)
	if 0 > ind || ind > max {
		return objects.NULL
	}
	return arr.Elements[ind]
}

func evalDeclarationStatement(node *ast.DeclarationStatement, environment *objects.Environment) objects.Object {
	include, ok := node.Declaration.(*ast.IncludeDeclaration)
	if !ok {
		return newError("unknown declaration type: %s", node.Declaration.String())
	}

	extEnv, err := ScriptLoader(include.Include.Value)
	if err != nil {
		return newError("unable preload external script:\n%s", err.Error())
	}

	environment.AddExternalEnvironment(include.Alias.Value, extEnv)

	return &objects.ExternalEnvironment{Environment: extEnv}
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
	result := make([]objects.Object, len(arguments))
	for i, argument := range arguments {
		evaluated := Eval(argument, environment)
		if isError(evaluated) {
			return []objects.Object{evaluated}
		}
		result[i] = evaluated
	}
	return result
}

func evalIdentifier(node *ast.Identifier, environment *objects.Environment) objects.Object {
	if val, ok := environment.Get(node.Value); ok {
		return val
	}

	if val, ok := environment.GetExternalEnvironment(node.Value); ok {
		return &objects.ExternalEnvironment{Environment: val}
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
	newEnv := objects.NewEnclosedEnvironment(environment)
	if isTruthy(condition) {
		return Eval(node.Consequence, newEnv)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, newEnv)
	}
	return objects.NULL
}

func evalForExpression(node *ast.ForExpression, environment *objects.Environment) objects.Object {
	newEnv := objects.NewEnclosedEnvironment(environment)

	if node.Initial != nil {
		Eval(node.Initial, newEnv)
	}

	var value objects.Object = objects.NULL
	for {
		if node.Condition != nil {
			cond := Eval(node.Condition, newEnv)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				return value
			}
		}

		value = Eval(node.Body, newEnv)
		if isError(value) {
			return value
		}

		if _, ok := value.(*objects.ReturnValue); ok {
			return value
		}

		if node.Complete != nil {
			compl := Eval(node.Complete, newEnv)
			if isError(value) {
				return compl
			}
		}
	}
}

func isTruthy(obj objects.Object) bool {
	switch obj {
	case objects.NULL, objects.FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(node *ast.InfixExpression, environment *objects.Environment) objects.Object {

	switch node.Operator {
	case ".":
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		return evalDottedExpression(left, node.Right, environment)
	case "=":
		return evalAssignExpression(node, environment)
	default:
		left := Eval(node.Left, environment)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, environment)
		if isError(right) {
			return right
		}
		return evalNativeInfixExpression(node.Operator, left, right)
	}
}

func evalAssignExpression(node *ast.InfixExpression, environment *objects.Environment) objects.Object {
	val := Eval(node.Right, environment)
	if isError(val) {
		return val
	}
	switch n := node.Left.(type) {
	case *ast.Identifier:
		value, ok := environment.Update(n.Value, val)
		if !ok {
			return newError("identifier not defined: %s", n.Value)
		}
		return value
	case *ast.IndexExpression:
		left := Eval(n.Left, environment)
		if isError(left) {
			return left
		}
		index := Eval(n.Index, environment)
		if isError(index) {
			return index
		}
		return evalAssignIndexExpression(left, index, val)
	case *ast.InfixExpression:
		return newError("unsupported multiple/inner/crosspackage assignments: %s", node.Right.TokenLiteral())
	default:
		return newError("unsupported assignment type receiver: %s", node.Right.TokenLiteral())
	}
}

func evalAssignIndexExpression(left objects.Object, index, value objects.Object) objects.Object {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return evalAssignArrayIndexExpression(left, index, value)
	case left.Type() == objects.HASH_OBJ:
		return evalAssignHashIndexExpression(left, index, value)
	default:
		return newError("index operator not supported for: %s", left.Type())
	}
}

func evalAssignHashIndexExpression(left objects.Object, index, value objects.Object) objects.Object {
	hash := left.(*objects.Hash)
	ind, ok := index.(objects.Hashable)
	if !ok {
		return newError("unusable as a hash key: %s", index.Type())
	}
	hash.Pairs[ind.HashKey()] = objects.HashPair{
		Key:   index,
		Value: value,
	}
	return value
}

func evalAssignArrayIndexExpression(left objects.Object, index, value objects.Object) objects.Object {
	arr := left.(*objects.Array)
	ind := index.(*objects.Integer).Value
	max := int64(len(arr.Elements) - 1)
	if 0 > ind || ind > max {
		return newError("index outbound: len=%d, ind=%d", max+1, ind)
	}
	arr.Elements[ind] = value
	return value
}

func evalDottedExpression(left objects.Object, right ast.Expression, environment *objects.Environment) objects.Object {
	switch leftExp := left.(type) {
	case *objects.ExternalEnvironment:
		switch n := right.(type) {
		case *ast.Identifier:
			return evalIdentifier(n, leftExp.Environment)
		case *ast.CallExpression:
			function := Eval(n.Function, leftExp.Environment)
			if isError(function) {
				return function
			}
			args := evalExpressions(n.Arguments, environment)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return applyFunction(function, args)
		case *ast.IndexExpression:
			left := Eval(n.Left, leftExp.Environment)
			if isError(left) {
				return left
			}
			index := Eval(n.Index, environment)
			if isError(index) {
				return index
			}
			return evalIndexExpression(left, index)
		default:
			return newError("unsupported reference call %s", right.TokenLiteral())
		}
	default:
		return newError("unsupported reference call on :%s", left.Type())
	}
}

func evalNativeInfixExpression(operator string, left objects.Object, right objects.Object) objects.Object {
	switch {
	case isNumbers(left.Type(), right.Type()):
		return evalNumberInfixExpression(operator, left, right)
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

func isNumbers(left, right objects.ObjectType) bool {
	return (left == objects.INTEGER_OBJ || left == objects.DOUBLE_OBJ) &&
		(right == objects.INTEGER_OBJ || right == objects.DOUBLE_OBJ)
}

func evalNumberInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := left.(objects.Arithmeticable)
	switch operator {
	case "+":
		return leftVal.Add(right)
	case "-":
		return leftVal.Sub(right)
	case "/":
		return leftVal.Div(right)
	case "*":
		return leftVal.Mul(right)
	}

	compLeft := left.(objects.Comparable)
	switch operator {
	case ">":
		return nativeBoolean(compLeft.Gt(right))
	case "<":
		return nativeBoolean(compLeft.Lt(right))
	case "==":
		return nativeBoolean(compLeft.Eq(right))
	case "!=":
		return nativeBoolean(compLeft.Neq(right))
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
	if value.Type() != objects.INTEGER_OBJ && value.Type() != objects.DOUBLE_OBJ {
		return newError("unknown operator: -%s", value.Type())
	}
	switch v := value.(type) {
	case *objects.Integer:
		return &objects.Integer{Value: -v.Value}
	case *objects.Double:
		return &objects.Double{Value: -v.Value}
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
