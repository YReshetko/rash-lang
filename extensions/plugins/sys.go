package main

import (
	"fmt"
	"time"
)

var SysPlugin = sysPlugin{}

const (
	pkg  = "sys"
	ver  = "0.0.1"
	desc = "provides system functions"
)

type sysPlugin struct{}

func (s sysPlugin) Eval(fnName string, args ...interface{}) ([]interface{}, error) {
	switch fnName {
	case "len":
		return s.length(args...)
	case "time":
		return s.time()
	default:
		return nil, fmt.Errorf("function %s not found in %s extension", fnName, pkg)
	}
}

func (s sysPlugin) Package() string {
	return pkg
}

func (s sysPlugin) Version() string {
	return ver
}

func (s sysPlugin) Description() string {
	return desc
}

func (s sysPlugin) length(args ...interface{}) ([]interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to `len`; got=%d, expected=1", len(args))
	}
	value := args[0]

	switch arg := value.(type) {
	case string:
		return []interface{}{len(arg)}, nil
	case []interface{}:
		return []interface{}{len(arg)}, nil
	case map[interface{}]interface{}:
		return []interface{}{len(arg)}, nil
	default:
		return nil, fmt.Errorf("unexpected value type %T to `len` function", value)
	}
}

func (s sysPlugin) time() ([]interface{}, error) {
	return []interface{}{time.Now().Format(time.RFC3339)}, nil
}
