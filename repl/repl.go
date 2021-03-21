package repl

import (
	"bufio"
	"fmt"
	"github.com/YReshetko/monkey-language/evaluator"
	"github.com/YReshetko/monkey-language/lexer"
	"github.com/YReshetko/monkey-language/objects"
	"github.com/YReshetko/monkey-language/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := objects.NewEnvironment()

	for {
		_, err := fmt.Fprint(out, PROMPT)
		if err != nil {
			return err
		}
		scanned := scanner.Scan()
		if !scanned {
			return nil
		}

		line := scanner.Text()
		if line == "exit" {
			return nil
		}

		l := lexer.New(line, "REPL")
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, s := range p.Errors() {
				_, err := fmt.Fprintf(out, "\t%s\n", s)
				if err != nil {
					return err
				}
			}
			continue
		}

		obj := evaluator.Eval(program, env)
		if obj != objects.NULL {
			_, err = fmt.Fprintf(out, "%s\n", obj.Inspect())
			if err != nil {
				return err
			}
		}
	}
}
