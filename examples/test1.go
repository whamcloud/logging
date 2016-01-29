package main

import (
	"flag"

	"github.intel.com/hpdd/debug"
)

var enableDebug bool

func init() {
	flag.BoolVar(&enableDebug, "debug", false, "enable debug logging")
}

func foo() {
	a := 123
	d := 123.123
	debug.Print("call from foo() ", a, d)

}

func main() {
	flag.Parse()
	if enableDebug {
		debug.Enable()
	}
	debug.Printf("inside %s", "main")
	foo()
}
