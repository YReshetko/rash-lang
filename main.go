package main

import (
	"fmt"
	"github.com/YReshetko/monkey-language/evaluator"
	"github.com/YReshetko/monkey-language/extensions"
	"github.com/YReshetko/monkey-language/repl"
	"log"
	"os"
	"os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	reg, err := extensionsRegistry()
	if err != nil {
		log.Fatal(err)
	}
	evaluator.InitRegistry(reg)

	fmt.Printf("Hello %s! Welcome in monkey language!\n", u.Username)
	fmt.Printf("Feel free to type any code here!\n")

	if err = repl.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good bye!")
}

func extensionsRegistry() (*extensions.Registry, error) {
	r := extensions.New()
	if err := r.Add("extensions/plugins/plugins.so", "SysPlugin"); err != nil {
		return nil, err
	}
	return r, nil
}
