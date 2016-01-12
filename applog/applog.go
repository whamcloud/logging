package applog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var taskSuffix = " ... "

type displayLevel int

func (d displayLevel) String() string {
	switch d {
	case DEBUG:
		return "DEBUG"
	case USER:
		return "USER"
	case WARN:
		return "WARN"
	case FAIL:
		return "FAIL"
	case SILENT:
		return "SILENT"
	default:
		return fmt.Sprintf("Unknown level: %d", d)
	}
}

const (
	// DEBUG shows all
	DEBUG displayLevel = iota
	// USER shows user-appropriate messages
	USER
	// WARN shows warnings
	WARN
	// FAIL is bad
	FAIL
	// SILENT shows nothing
	SILENT
)

type loggedWriter struct {
	prefix string
	logger *AppLogger
}

func (w *loggedWriter) Write(data []byte) (int, error) {
	w.logger.Debug(fmt.Sprintf("%s %s", w.prefix, data))
	return len(data), nil
}

func newLoggedWriter(prefix string, logger *AppLogger) *loggedWriter {
	return &loggedWriter{
		prefix: prefix,
		logger: logger,
	}
}

// OptSetter sets logger options
type OptSetter func(*AppLogger)

// JournalFile configures the logger's journaler to use the specified file
func JournalFile(w io.Writer) OptSetter {
	return func(l *AppLogger) {
		l.journal = log.New(w, "", log.LstdFlags)
	}
}

// DisplayLevel sets the logger's display level
func DisplayLevel(d displayLevel) OptSetter {
	return func(l *AppLogger) {
		l.Level = d
	}
}

// New returns a new AppLogger
func New(options ...OptSetter) *AppLogger {
	logger := &AppLogger{
		spinner: spinner.New(spinner.CharSets[9], 100*time.Millisecond),
		out:     os.Stdout,
		err:     os.Stderr,
		Level:   USER,
		journal: log.New(ioutil.Discard, "", log.LstdFlags),
	}

	for _, option := range options {
		option(logger)
	}

	return logger
}

// AppLogger is a logger with methods for displaying entries to the user
// after recording them to a journal.
type AppLogger struct {
	Level displayLevel

	spinner     *spinner.Spinner
	out         io.Writer
	err         io.Writer
	lastEntry   string
	currentTask string
	journal     *log.Logger
}

// Out returns an io.Writer that can capture output
func (l *AppLogger) Out() io.Writer {
	return newLoggedWriter(">", l)
}

// Err returns an io.Writer that can capture output
func (l *AppLogger) Err() io.Writer {
	return newLoggedWriter("!", l)
}

// DisplayLevel sets the logger's display level
func (l *AppLogger) DisplayLevel(level displayLevel) {
	DisplayLevel(level)(l)
}

// JournalFile configures the logger's journaler to use the specified file
func (l *AppLogger) JournalFile(w io.Writer) {
	JournalFile(w)(l)
}

func (l *AppLogger) recordEntry(level displayLevel, v ...interface{}) {
	if len(v) == 0 {
		return
	}

	switch arg := v[0].(type) {
	case error:
		l.lastEntry = fmt.Sprintf("ERROR: %s", arg)
	case string:
		if len(v) > 1 {
			l.lastEntry = fmt.Sprintf(arg, v[1:]...)
		} else {
			l.lastEntry = fmt.Sprint(arg)
		}
	default:
		panic(fmt.Sprintf("Unhandled entry: %v", v))
	}
	l.journal.Printf("%s: %s", level, l.lastEntry)
}

// Debug logs the entry and prints to stdout if level <= DEBUG
func (l *AppLogger) Debug(v ...interface{}) {
	l.recordEntry(DEBUG, v...)

	if l.Level <= DEBUG {
		fmt.Fprintf(l.out, "%s: %s\n", DEBUG, l.lastEntry)
	}
}

// User logs the entry and prints to stdout if level <= USER
func (l *AppLogger) User(v ...interface{}) {
	l.recordEntry(USER, v...)

	if l.Level <= USER {
		fmt.Fprintln(l.out, l.lastEntry)
	}
}

// StartTask logs the entry at USER level and displays a spinner
// for long-running tasks
func (l *AppLogger) StartTask(v ...interface{}) {
	// Allow new tasks to display completion for previous tasks.
	if l.currentTask != "" {
		l.CompleteTask()
	}

	l.recordEntry(USER, v...)

	if l.Level == USER {
		l.currentTask = l.lastEntry
		l.spinner.Prefix = l.currentTask + taskSuffix
		l.spinner.Restart()
	}
}

// CompleteTask stops the spinner and prints a newline
func (l *AppLogger) CompleteTask(v ...interface{}) {
	l.spinner.Stop()

	if len(v) == 0 {
		l.recordEntry(USER, l.currentTask+taskSuffix+"Done.")
	} else {
		if fmtStr, ok := v[0].(string); ok {
			var newArgs []interface{}
			newArgs = append(newArgs, l.currentTask+taskSuffix+fmtStr)
			newArgs = append(newArgs, v[1:]...)
			l.recordEntry(USER, newArgs...)
		} else {
			l.recordEntry(USER, v...)
		}
	}

	if l.currentTask != "" && l.Level == USER {
		fmt.Fprintln(l.out, l.lastEntry)
		l.currentTask = ""
	}
}

// Warn logs the entry and prints to stderr if level <= WARN
func (l *AppLogger) Warn(v ...interface{}) {
	l.recordEntry(WARN, v...)

	l.spinner.Stop()
	l.currentTask = ""
	if l.Level <= WARN {
		fmt.Fprintf(l.err, "%s: %s\n", WARN, l.lastEntry)
	}
}

// Fail logs the entry and prints to stderr if level <= FAIL
func (l *AppLogger) Fail(v ...interface{}) {
	l.recordEntry(FAIL, v...)

	l.spinner.Stop()
	l.currentTask = ""
	if l.Level <= FAIL {
		fmt.Fprintln(l.err, l.lastEntry)
	}
	os.Exit(1)
}

var std = New()

// StandardLogger returns the standard logger configured by the library
func StandardLogger() *AppLogger {
	return std
}

// SetStandard sets the standard logger to the supplied logger
func SetStandard(l *AppLogger) {
	std = l
}

// Debug logs the entry and prints to stdout if level <= DEBUG
func Debug(v ...interface{}) {
	std.Debug(v...)
}

// User logs the entry and prints to stdout if level <= USER
func User(v ...interface{}) {
	std.User(v...)
}

// Warn logs the entry and prints to stderr if level <= WARN
func Warn(v ...interface{}) {
	std.Warn(v...)
}

// Fail logs the entry and prints to stderr if level <= FAIL
func Fail(v ...interface{}) {
	std.Fail(v...)
}
