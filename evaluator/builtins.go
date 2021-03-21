package evaluator

import (
	"github.com/YReshetko/monkey-language/extensions"
	"github.com/YReshetko/monkey-language/objects"
)

var registry *extensions.Registry

func InitRegistry(r *extensions.Registry) {
	registry = r
}

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
			return retVal(returnVal[0])
		},
	},
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
	default:
		return nil
	}
}
