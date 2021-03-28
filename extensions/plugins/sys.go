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
type Callback func(args ...interface{}) ([]interface{}, error)

func (s sysPlugin) Eval(fnName string, args ...interface{}) ([]interface{}, error) {
	switch fnName {
	case "len":
		return s.length(args...)
	case "time":
		return s.time()
	case "print":
		return s.print(args...)
	default:
		return nil, fmt.Errorf("function %s not found in %s extension", fnName, pkg)
	}
}

func (s sysPlugin) Call(fnName string, callback func(args ...interface{}) ([]interface{}, error), args ...interface{}) ([]interface{}, error) {
	switch fnName {
	case "tick":
		return s.tick(callback, args...)
	default:
		return nil, fmt.Errorf("callback function %s not found in %s extension", fnName, pkg)
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

func (s sysPlugin) print(args ...interface{}) ([]interface{}, error) {
	fmt.Println(args...)
	return nil, nil
}

func (s sysPlugin) tick(callback Callback, args ...interface{}) ([]interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to `tick`; got=%d, expected=1", len(args))
	}
	duration, ok := args[0].(int64)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type sent to `tick`; got=%d, expected=int", args[0])
	}

	ticker := time.NewTicker(time.Second * time.Duration(duration))
	go func() {
		for {
			select {
			case x := <-ticker.C:
				_, err := callback(x.Format(time.RFC3339))
				if err != nil {
					ticker.Stop()
					return
				}
			}
		}
	}()

	return nil, nil
}
