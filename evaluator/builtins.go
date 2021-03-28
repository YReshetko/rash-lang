package evaluator

import (
	"errors"
	"github.com/YReshetko/rash-lang/ast"
	"github.com/YReshetko/rash-lang/extensions"
	"github.com/YReshetko/rash-lang/objects"
)

var registry *extensions.Registry

func InitRegistry(r *extensions.Registry) {
	registry = r
}

type Evaluator func(node ast.Node, environment *objects.Environment) objects.Object
var Evaluate Evaluator

var builtins = map[string]*objects.Builtin{
	"eval": { // eval function expects at leas two arguments plugin_name and called_function, all others will be passed to plugin as function call
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments to `eval`; got=%d, expected>=%d", len(args), 2)
			}
			pkgName, ok := args[0].(*objects.String)
			if !ok {
				return newError("`eval` expects string as first argument, but got %s", args[0].Type())
			}
			fnName, ok := args[1].(*objects.String)
			if !ok {
				return newError("`eval` expects string as second argument, but got %s", args[1].Type())
			}

			inArgs := []interface{}{}
			for i := 2; i < len(args); i++ {
				inArgs = append(inArgs, getValue(args[i]))
			}

			returnVal, err := registry.Eval(pkgName.Value, fnName.Value, inArgs...)

			if err != nil {
				return newError("plugin `%s` err: %v", pkgName.Value, err)
			}
			if len(returnVal) == 0 {
				return objects.NULL
			}
			// Suppose the fires value has meaning
			// TODO make array/map mappable to `rash` array
			return retVal(returnVal[0])
		},
	},
	"call": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) < 3 {
				return newError("wrong number of arguments to `call`; got=%d, expected>=%d", len(args), 3)
			}
			pkgName, ok := args[0].(*objects.String)
			if !ok {
				return newError("`call` expects string as first argument, but got %s", args[0].Type())
			}
			fnName, ok := args[1].(*objects.String)
			if !ok {
				return newError("`call` expects string as second argument, but got %s", args[1].Type())
			}
			fn, ok := args[2].(*objects.Function)

			inArgs := []interface{}{}
			for i := 3; i < len(args); i++ {
				inArgs = append(inArgs, getValue(args[i]))
			}

			retValue, err := registry.Call(pkgName.Value, fnName.Value, NewCallback(fn), inArgs...)
			if err != nil {
				return newError("plugin `%s` err: %v", pkgName.Value, err)
			}
			if len(retValue) == 0 {
				return objects.NULL
			}
			// TODO make array/map mappable to `rash` array
			return retVal(retValue[0])
		},
	},
}

func NewCallback(fn *objects.Function) func(args ...interface{}) ([]interface{}, error) {
	return func(args ...interface{}) ([]interface{}, error) {
		prepArgs := make([]objects.Object, len(args))
		for i, v := range args {
			prepArgs[i] = retVal(v)
		}

		if len(prepArgs) != len(fn.Parameters) {
			return nil, errors.New("unexpected number of arguments")
		}
		extendedEnv := extendFunctionEnvironment(fn, prepArgs)
		evaluated := Evaluate(fn.Body, extendedEnv)

		outValues := []interface{}{}

		if evaluated != objects.NULL{
			outValues = []interface{}{getValue(evaluated)}
		}

		return outValues, nil
	}
}

func retVal(val interface{}) objects.Object {
	switch v := val.(type) {
	case int:
		return &objects.Integer{Value: int64(v)}
	case int64:
		return &objects.Integer{Value: v}
	case int32:
		return &objects.Integer{Value: int64(v)}
	case int8:
		return &objects.Integer{Value: int64(v)}
	case int16:
		return &objects.Integer{Value: int64(v)}
	case byte:
		return &objects.Integer{Value: int64(v)}
	case uint64:
		return &objects.Integer{Value: int64(v)}
	case uint32:
		return &objects.Integer{Value: int64(v)}
	case uint:
		return &objects.Integer{Value: int64(v)}
	case *int:
		return &objects.Integer{Value: int64(*v)}
	case *int64:
		return &objects.Integer{Value: *v}
	case *int32:
		return &objects.Integer{Value: int64(*v)}
	case *int8:
		return &objects.Integer{Value: int64(*v)}
	case *int16:
		return &objects.Integer{Value: int64(*v)}
	case *byte:
		return &objects.Integer{Value: int64(*v)}
	case *uint64:
		return &objects.Integer{Value: int64(*v)}
	case *uint32:
		return &objects.Integer{Value: int64(*v)}
	case *uint:
		return &objects.Integer{Value: int64(*v)}
	case string:
		return &objects.String{Value: v}
	case *string:
		return &objects.String{Value: *v}
	case bool:
		return nativeBoolean(v)
	case *bool:
		return nativeBoolean(*v)
	default:
		return objects.NULL
	}
}

func getValue(object objects.Object) interface{} {
	switch obj := object.(type) {
	case *objects.String:
		return obj.Value
	case *objects.Integer:
		return obj.Value
	case *objects.Boolean:
		return obj.Value
	case *objects.Array:
		arr := make([]interface{}, len(obj.Elements))
		for i, element := range obj.Elements {
			arr[i] = getValue(element)
		}
		return arr
	case *objects.Hash:
		m := map[interface{}]interface{}{}
		for _, v := range obj.Pairs{
			m[getValue(v.Key)] = getValue(v.Value)
		}
		return m
	default:
		return nil
	}
}
