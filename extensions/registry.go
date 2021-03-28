package extensions

import (
	"errors"
	"fmt"
	"plugin"
)

type Registry struct {
	plugins map[string]Plugin
}

func New() *Registry {
	return &Registry{plugins: map[string]Plugin{}}
}

func (r *Registry) Add(file, symbol string) error {
	p, err := plugin.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open plugin file: %v", err)
	}

	sym, err := p.Lookup(symbol)
	if err != nil {
		return fmt.Errorf("unable to find exported symbol: %v", err)
	}

	plug, ok := sym.(Plugin)
	if !ok {
		return errors.New("exported symbol doesn't match plugin interface")
	}

	// TODO version validation
	r.plugins[plug.Package()] = plug
	return nil
}

func (r *Registry) Eval(pkgName, fnName string, args ...interface{}) ([]interface{}, error) {
	plug, ok := r.plugins[pkgName]
	if !ok {
		return nil, fmt.Errorf("package %s not found in extensions", pkgName)
	}

	return plug.Eval(fnName, args...)
}


func (r *Registry) Call(pkgName, fnName string, fn func(args ...interface{}) ([]interface{}, error), args ...interface{}) ([]interface{}, error) {
	plug, ok := r.plugins[pkgName]
	if !ok {
		return nil, fmt.Errorf("package %s not found in extensions", pkgName)
	}
	return plug.Call(fnName, fn, args...)
}
