package main

import (
	"flag"
	"log"
	"os"

	"github.intel.com/hpdd/debug"
)

var enableDebug bool

var dbg *debug.Debugger

func init() {
	flag.BoolVar(&enableDebug, "debug", false, "enable debug logging")
	// Localized debugger
	dbg = debug.New(log.New(os.Stderr, "", log.LstdFlags|log.Llongfile))
}

func foo() {
	a := 123
	d := 123.123
	debug.Print("call from foo() ", a, d)
	debug.Printf("call from foo() %v %v", a, d)

}

func bar() {
	a := 123
	dbg.Print("local debugger")
	dbg.Printf("format: %d", a)
}

func main() {
	flag.Parse()
	if enableDebug {
		debug.Enable()
		dbg.Enable()
	}
	debug.Printf("inside %s", "main")
	foo()
	bar()
}
