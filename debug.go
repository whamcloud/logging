package debug

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type (
	// Config holds configuration values for the package
	Debugger struct {
		mu      sync.Mutex
		log     *log.Logger
		enabled bool
	}
)

var std *Debugger

func init() {
	std = New(log.New(os.Stdout, "DEBUG ", log.Lmicroseconds|log.Lshortfile))
}

func New(log *log.Logger) *Debugger {
	return &Debugger{log: log}
}

func (d *Debugger) SetLogger(log *log.Logger) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.log = log
}

func (d *Debugger) Enabled() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.enabled
}

func (d *Debugger) Enable() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = true
}

func (d *Debugger) Disable() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.enabled = false
}

func (d *Debugger) Output(skip int, v ...interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.enabled {
		return
	}
	d.log.Output(skip, fmt.Sprint(v...))
}

func (d *Debugger) Outputf(skip int, f string, v ...interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.enabled {
		return
	}
	d.log.Output(skip, fmt.Sprintf(f, v...))
}

func (d *Debugger) Printf(f string, v ...interface{}) {
	d.Outputf(3, f, v...)
}

func (d *Debugger) Print(v ...interface{}) {
	d.Output(3, v...)
}

func (d *Debugger) Assertf(expr bool, f string, v ...interface{}) {
	if !d.Enabled() {
		return
	}
	if !expr {
		msg := fmt.Sprintf("ASSERTION FAILED: "+f, v...)
		d.Output(3, msg)
		panic(msg)
	}
}

func (d *Debugger) Assert(expr bool, v ...interface{}) {
	if !d.Enabled() {
		return
	}
	if !expr {
		msg := fmt.Sprintf("ASSERTION FAILED: %s", fmt.Sprint(v...))
		d.Output(3, msg)
		panic(msg)
	}
}

// SetDebugger replaces the debuggeruration
func SetLogger(log *log.Logger) {
	std.SetLogger(log)
}

// EnableDebug enables debug logging
func Enable() {
	std.Enable()
}

// DisableDebug disables debug logging
func Disable() {
	std.Disable()
}

// Printf prints message if debug logging is enabled.
func Printf(f string, v ...interface{}) {
	std.Outputf(3, f, v...)
}

// Print prints arguments if debug logging is enabled.
func Print(v ...interface{}) {
	std.Output(3, v...)
}

// Assertf will panic if expression is not true, but only if debugging is enabled
func Assertf(expr bool, f string, v ...interface{}) {
	if !std.Enabled() {
		return
	}
	if !expr {
		msg := fmt.Sprintf("ASSERTION FAILED: "+f, v...)
		std.Output(3, msg)
		panic(msg)
	}
}

// Assert will panic if expression is not true, but only if debugging is enabled
func Assert(expr bool, v ...interface{}) {
	if !std.Enabled() {
		return
	}
	if !expr {
		msg := fmt.Sprintf("ASSERTION FAILED: %s", fmt.Sprint(v...))
		std.Output(3, msg)
		panic(msg)
	}
}
