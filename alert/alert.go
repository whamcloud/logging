package alert

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

type (
	// Logger defines an interface for an alert logger
	Logger interface {
		SetFlags(int)
		SetOutput(io.Writer)
		Output(int, string)

		Warn(...interface{})
		Warnf(string, ...interface{})

		Fatal(...interface{})
		Fatalf(string, ...interface{})
	}

	// StdErrLogger is a logger which writes output to os.Stderr
	StdErrLogger struct {
		log *log.Logger
	}
)

var std Logger

// Provide as much information as possible about where the message originated,
// as this package should usually only be involved where there is a failure.
const logFlags = log.Ldate | log.Ltime | log.LUTC | log.Llongfile

func init() {
	std = NewStdErrLogger()
}

// NewStdErrLogger returns a *StdErrLogger
func NewStdErrLogger() *StdErrLogger {
	return &StdErrLogger{
		log: log.New(os.Stderr, "ALERT ", logFlags),
	}
}

// SetFlags sets the output flags for the embedded logger
func (l *StdErrLogger) SetFlags(flags int) {
	l.log.SetFlags(flags)
}

// SetOutput updates the embedded logger's output
func (l *StdErrLogger) SetOutput(out io.Writer) {
	l.log.SetOutput(out)
}

// Output writes the output for a logging event
func (l *StdErrLogger) Output(skip int, s string) {
	l.log.Output(skip, s)
}

// Warn outputs a log message from the arguments
func (l *StdErrLogger) Warn(v ...interface{}) {
	l.Output(3, fmt.Sprint(v...))
}

// Warnf outputs a formatted log message from the arguments
func (l *StdErrLogger) Warnf(f string, v ...interface{}) {
	l.Output(3, fmt.Sprintf(f, v...))
}

// Fatal outputs a log message from the arguments, then exits
func (l *StdErrLogger) Fatal(v ...interface{}) {
	l.Output(3, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf outputs a formatted log message from the arguments, then exits
func (l *StdErrLogger) Fatalf(f string, v ...interface{}) {
	l.Output(3, fmt.Sprintf(f, v...))
	os.Exit(1)
}

// package-level functions follow

// SetOutput configures the output writer for the logger
func SetOutput(out io.Writer) {
	std.SetOutput(out)
}

// Warn outputs a log message from the arguments
func Warn(v ...interface{}) {
	std.Output(3, fmt.Sprint(v...))
}

// Warnf outputs a formatted log message from the arguments
func Warnf(f string, v ...interface{}) {
	std.Output(3, fmt.Sprintf(f, v...))
}

// Fatal outputs a log message from the arguments, then exits
func Fatal(v ...interface{}) {
	std.Output(3, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf outputs a formatted log message from the arguments, then exits
func Fatalf(f string, v ...interface{}) {
	std.Output(3, fmt.Sprintf(f, v...))
	os.Exit(1)
}

// Abort prints error trace and exits
func Abort(err error) {
	var b bytes.Buffer

	// We don't need to see where the abort was called, so we remove
	// this flag before logging and exiting.
	std.SetFlags(logFlags &^ log.Llongfile)

	errors.Fprint(&b, err)
	std.Output(3, "Aborting program execution due to error(s):\n"+b.String())
	os.Exit(1)
}
