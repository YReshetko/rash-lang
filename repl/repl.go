package repl

import (
	"bufio"
	"fmt"
	"github.com/YReshetko/rash-lang/evaluator"
	"github.com/YReshetko/rash-lang/lexer"
	"github.com/YReshetko/rash-lang/objects"
	"github.com/YReshetko/rash-lang/parser"
	"io"
)

const PROMPT = ">> "

const initial = `
# http "http.rs";
# sys "imports.rs";
let server1 = http.new_server("4000");
server["register"]("GET", "/hello", fn(){return "Hello world"});
server["start"]();

`

/*

*/

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := objects.NewEnvironment()

	eval(initial, env, out)

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

		if len(line) == 0 {
			continue
		}

		eval(line, env, out)
	}
}

func eval(input string, environment *objects.Environment, out io.Writer) {
	l := lexer.New(initial, "REPL")
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, s := range p.Errors() {
			_, err := fmt.Fprintf(out, "\t%s\n", s)
			if err != nil {
				return
			}
		}
		return
	}

	obj := evaluator.Eval(program, environment)
	if obj != objects.NULL {
		_, err := fmt.Fprintf(out, "%s\n", obj.Inspect())
		if err != nil {
			return
		}
	}

}
