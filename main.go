package main

import (
	"fmt"
	"github.com/YReshetko/rash-lang/evaluator"
	"github.com/YReshetko/rash-lang/extensions"
	"github.com/YReshetko/rash-lang/loaders"
	"github.com/YReshetko/rash-lang/repl"
	"log"
	"os"
	"os/user"
)

const banner = `
@@@@@@@    @@@@@@    @@@@@@   @@@  @@@  @@@@@@@@   @@@@@@   @@@  @@@   @@@@@@   
@@@@@@@@  @@@@@@@@  @@@@@@@   @@@  @@@  @@@@@@@@  @@@@@@@   @@@  @@@  @@@@@@@@  
@@!  @@@  @@!  @@@  !@@       @@!  @@@  @@!       !@@       @@!  !@@  @@!  @@@  
!@!  @!@  !@!  @!@  !@!       !@!  @!@  !@!       !@!       !@!  @!!  !@!  @!@  
@!@!!@!   @!@!@!@!  !!@@!!    @!@!@!@!  @!!!:!    !!@@!!    @!@@!@!   @!@!@!@!  
!!@!@!    !!!@!!!!   !!@!!!   !!!@!!!!  !!!!!:     !!@!!!   !!@!!!    !!!@!!!!  
!!: :!!   !!:  !!!       !:!  !!:  !!!  !!:            !:!  !!: :!!   !!:  !!!  
:!:  !:!  :!:  !:!      !:!   :!:  !:!  :!:           !:!   :!:  !:!  :!:  !:!  
::   :::  ::   :::  :::: ::   ::   :::   :: ::::  :::: ::    ::  :::  ::   :::  
 :   : :   :   : :  :: : :     :   : :  : :: ::   :: : :     :   :::   :   : :  
`

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	reg, err := extensionsRegistry()
	if err != nil {
		log.Fatal(err)
	}

	evaluator.ScriptLoader = loaders.ScriptLoader
	evaluator.Evaluate = evaluator.Eval
	evaluator.InitRegistry(reg)

	fmt.Println(banner)
	fmt.Printf("Hello %s! Welcome in `rasheska` script language!\n", u.Username)
	fmt.Printf("Let's start fun!\n")

	if err = repl.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Good bye!... rasheska will miss you")
}

func extensionsRegistry() (*extensions.Registry, error) {
	r := extensions.New()
	if err := r.Add("bin/sys.so", "SysPlugin"); err != nil {
		return nil, err
	}
	return r, nil
}
