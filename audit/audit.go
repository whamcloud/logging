package audit

import (
	"fmt"
	"log"
	"os"
)

type (
	// Logger defines an interface for an audit logger
	Logger interface {
		Output(int, string)

		Log(...interface{})
		Logf(string, ...interface{})
	}

	// StdOutLogger is a logger which writes output to os.Stdout
	StdOutLogger struct {
		log *log.Logger
	}
)

var std Logger

const logFlags = log.LstdFlags | log.LUTC

func init() {
	std = NewStdOutLogger()
}

// NewStdOutLogger returns a *StdOutLogger
func NewStdOutLogger() *StdOutLogger {
	return &StdOutLogger{
		log: log.New(os.Stdout, "", logFlags),
	}
}

// Output writes the output for a logging event
func (l *StdOutLogger) Output(skip int, s string) {
	l.log.Output(skip, s)
}

// Log outputs a log message from the arguments
func (l *StdOutLogger) Log(v ...interface{}) {
	l.Output(3, fmt.Sprint(v...))
}

// Logf outputs a formatted log message from the arguments
func (l *StdOutLogger) Logf(f string, v ...interface{}) {
	l.Output(3, fmt.Sprintf(f, v...))
}

// package-level functions follow

// Log outputs a log message from the arguments
func Log(v ...interface{}) {
	std.Output(3, fmt.Sprint(v...))
}

// Logf outputs a formatted log message from the arguments
func Logf(f string, v ...interface{}) {
	std.Output(3, fmt.Sprintf(f, v...))
}
