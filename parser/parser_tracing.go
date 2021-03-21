package parser

import "fmt"

var (
	isTrace = false
	tab     = ""
)

func trace(method string) string {
	if !isTrace {
		return ""
	}
	inc()
	fmt.Printf("%sBEGIN: %s;\n", tab, method)
	return method
}

func untrace(method string) {
	if !isTrace {
		return
	}
	fmt.Printf("%sEND  : %s;\n", tab, method)
	dec()
}

func inc() {
	tab += "\t"
}

func dec() {
	if len(tab) != 0 {
		tab = tab[1:]
	}
}
