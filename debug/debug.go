package debug

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

type (
	// Debugger wraps a *log.Logger with some configuration and
	// convenience methods
	Debugger struct {
		sync.Mutex
		log     *log.Logger
		enabled bool
	}

	// Flag allows the flag package to enable debugging
	Flag bool
)

var std *Debugger

func init() {
	std = New(log.New(os.Stdout, "DEBUG ", log.Lmicroseconds|log.Lshortfile))
}

// FlagVar returns a tuple of parameters suitable for flag.Var()
func FlagVar() (*Flag, string, string) {
	f := Flag(false)
	return &f, "debug", "enable debug output"
}

// IsBoolFlag satisfies the flag.boolFlag interface
func (f *Flag) IsBoolFlag() bool {
	return true
}

func (f *Flag) String() string {
	return fmt.Sprintf("%v", *f)
}

// Set satisfies the flag.Value interface
func (f *Flag) Set(value string) error {
	b, err := strconv.ParseBool(value)
	if err == nil {
		std.enabled = b
		f = (*Flag)(&b)
	}
	return err
}

// New wraps a *log.Logger with a *debug.Debugger
func New(log *log.Logger) *Debugger {
	return &Debugger{log: log}
}

// SetLogger accepts a new *log.Logger to wrap
func (d *Debugger) SetLogger(log *log.Logger) {
	d.Lock()
	defer d.Unlock()
	d.log = log
}

// Enabled indicates whether or not debugging is enabled
func (d *Debugger) Enabled() bool {
	d.Lock()
	defer d.Unlock()
	return d.enabled
}

// Enable turns on debug logging
func (d *Debugger) Enable() {
	d.Lock()
	defer d.Unlock()
	d.enabled = true
}

// Disable turns off debug logging
func (d *Debugger) Disable() {
	d.Lock()
	defer d.Unlock()
	d.enabled = false
}

// Output writes the output for a logging event
func (d *Debugger) Output(skip int, s string) {
	if !d.Enabled() {
		return
	}
	d.log.Output(skip, s)
}

// Printf outputs formatted arguments
func (d *Debugger) Printf(f string, v ...interface{}) {
	d.Output(3, fmt.Sprintf(f, v...))
}

// Print outputs the arguments
func (d *Debugger) Print(v ...interface{}) {
	d.Output(3, fmt.Sprint(v...))
}

// Assertf accepts a boolean expression and formatted arguments, which
// if the expression is false, will be printed before panicing.
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

// Assert accepts a boolean expression and arguments, which if the
// expression is false, will be printed before panicing.
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

// SetLogger replaces the wrapped *log.Logger
func SetLogger(log *log.Logger) {
	std.SetLogger(log)
}

// Enable enables debug logging
func Enable() {
	std.Enable()
}

// Disable disables debug logging
func Disable() {
	std.Disable()
}

// Enabled returns a bool indicating whether or not debugging is enabled
func Enabled() bool {
	return std.Enabled()
}

// Printf prints message if debug logging is enabled.
func Printf(f string, v ...interface{}) {
	std.Output(3, fmt.Sprintf(f, v...))
}

// Print prints arguments if debug logging is enabled.
func Print(v ...interface{}) {
	std.Output(3, fmt.Sprint(v...))
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
