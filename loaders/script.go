package loaders

import (
	"fmt"
	"github.com/YReshetko/rash-lang/evaluator"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/objects"
	"github.com/YReshetko/rash-lang/parser"
	"io/ioutil"
	"strings"
)

func ScriptLoader(path string) (*objects.Environment, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load included script %s due to %v", path, err)
	}

	l := lexer.New(string(src), path)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("unable to evaluate included script %s due to:\n %s", path, strings.Join(p.Errors(), ";\n"))
	}

	externalEnv := objects.NewEnvironment()

	obj := evaluator.Eval(program, externalEnv)
	if obj.Type() == objects.ERROR_OBJ {
		errObj := obj.(*objects.Error)
		return nil, fmt.Errorf("%s\nStackTrace:\n%s", errObj.Inspect(), strings.Join(errObj.Stack, ";\n"))
	}

	return externalEnv, nil
}
