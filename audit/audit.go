package audit

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.intel.com/hpdd/logging/external"
)

type (
	config struct {
		sync.Mutex
		out       io.Writer
		externals []*external.Writer
	}

	// Logger defines an interface for an audit logger
	Logger interface {
		SetOutput(io.Writer)
		Output(int, string)

		Log(...interface{})
		Logf(string, ...interface{})
	}

	// StdOutLogger is a logger which writes output to os.Stdout
	StdOutLogger struct {
		log *log.Logger
	}

	// ExternalWriter is an optionally-prefixed writer for
	// 3rd-party logging packages.
	ExternalWriter struct {
		log *log.Logger
	}
)

var (
	cfg *config
	std Logger
)

const logFlags = log.LstdFlags | log.LUTC

func init() {
	cfg = &config{
		out: os.Stdout,
	}

	std = NewStdOutLogger()
}

// NewStdOutLogger returns a *StdOutLogger
func NewStdOutLogger() *StdOutLogger {
	cfg.Lock()
	defer cfg.Unlock()

	return &StdOutLogger{
		log: log.New(cfg.out, "", logFlags),
	}
}

// SetOutput updates the embedded logger's output
func (l *StdOutLogger) SetOutput(out io.Writer) {
	l.log.SetOutput(out)
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

// Writer returns a new *external.Writer suitable for injection into
// 3rd-party logging packages.
func Writer() *external.Writer {
	cfg.Lock()
	defer cfg.Unlock()

	w := external.NewWriter(cfg.out)
	cfg.externals = append(cfg.externals, w)
	return w
}

// Log outputs a log message from the arguments
func Log(v ...interface{}) {
	std.Output(3, fmt.Sprint(v...))
}

// Logf outputs a formatted log message from the arguments
func Logf(f string, v ...interface{}) {
	std.Output(3, fmt.Sprintf(f, v...))
}

// SetOutput updates the io.Writer for the package as well as any external
// writers created by the package
func SetOutput(out io.Writer) {
	cfg.Lock()
	defer cfg.Unlock()

	cfg.out = out
	std.SetOutput(cfg.out)
	for _, writer := range cfg.externals {
		writer.SetOutput(cfg.out)
	}
}
