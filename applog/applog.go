package applog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

var (
	taskSuffix = " ... "
	std        *AppLogger
)

func init() {
	std = New()
}

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

// LoggedWriter implements io.Writer and is used to redirect logging from
// 3rd-party libraries to this library.
type LoggedWriter struct {
	level  displayLevel
	prefix string
	logger *AppLogger
}

// Write logs the data at the specified loglevel
func (w *LoggedWriter) Write(data []byte) (int, error) {
	msg := string(data)
	if len(w.prefix) > 0 {
		msg = fmt.Sprintf("%s %s", w.prefix, data)
	}
	w.logger.logAt(w.level, msg)

	return len(data), nil
}

// Prefix optionally sets the LoggedWriter prefix
func (w *LoggedWriter) Prefix(prefix string) *LoggedWriter {
	w.prefix = prefix
	return w
}

// Level optionally sets the LoggedWriter log level
func (w *LoggedWriter) Level(level displayLevel) *LoggedWriter {
	w.level = level
	return w
}

// OptSetter sets logger options
type OptSetter func(*AppLogger)

// JournalFile configures the logger's journaler
func JournalFile(w interface{}) OptSetter {
	var writer io.Writer
	switch w := w.(type) {
	case io.Writer:
		writer = w
	case string:
		switch strings.ToLower(w) {
		case "":
			writer = ioutil.Discard
		case "stderr":
			writer = os.Stderr
		case "stdout":
			writer = os.Stdout
		default:
			file, err := os.OpenFile(w, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
			if err != nil {
				panic(fmt.Errorf("applog.JournalFile() failed to open %s: %s", w, err))
			}
			writer = file
		}
	default:
		panic(fmt.Errorf("applog.JournalFile() called with invalid argument: %s", w))
	}

	return func(l *AppLogger) {
		l.journal = log.New(writer, "", log.LstdFlags)
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

func (l *AppLogger) logAt(level displayLevel, msg string) {
	switch level {
	case SILENT:
		return
	case USER:
		l.User(msg)
	case WARN:
		l.Warn(msg)
	case FAIL:
		l.Fail(msg)
	default:
		l.Debug(msg)
	}
}

// Out returns an io.Writer that can capture prefixed stdout at DEBUG level
func (l *AppLogger) Out() io.Writer {
	return l.Writer().Level(DEBUG).Prefix(">")
}

// Err returns an io.Writer that can capture prefixed stderr at DEBUG level
func (l *AppLogger) Err() io.Writer {
	return l.Writer().Level(DEBUG).Prefix("!")
}

// Writer returns an io.Writer for injecting our logging into third-party
// libraries
func (l *AppLogger) Writer() *LoggedWriter {
	return &LoggedWriter{
		level:  DEBUG,
		prefix: "",
		logger: l,
	}
}

// DisplayLevel sets the logger's display level
func (l *AppLogger) DisplayLevel(level displayLevel) {
	DisplayLevel(level)(l)
}

// JournalFile configures the logger's journaler
func (l *AppLogger) JournalFile(w interface{}) {
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
		fmt.Fprintf(l.out, "%s: %s", DEBUG, l.lastEntry)
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
		fmt.Fprintf(l.err, "%s: %s", WARN, l.lastEntry)
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

// StandardLogger returns the standard logger configured by the library
func StandardLogger() *AppLogger {
	return std
}

// SetStandard sets the standard logger to the supplied logger
func SetStandard(l *AppLogger) {
	std = l
}

// SetJournal sets the standard logger's journal writer
func SetJournal(w interface{}) {
	JournalFile(w)(std)
}

// SetLevel sets the standard logger's display level
func SetLevel(d displayLevel) {
	DisplayLevel(d)(std)
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

// StartTask logs the entry at USER level and displays a spinner
// for long-running tasks
func StartTask(v ...interface{}) {
	std.StartTask(v...)
}

// CompleteTask stops the spinner and prints a newline
func CompleteTask(v ...interface{}) {
	std.CompleteTask(v...)
}

// Out returns an io.Writer that can capture prefixed stdout at DEBUG level
func Out() io.Writer {
	return std.Out()
}

// Err returns an io.Writer that can capture prefixed stderr at DEBUG level
func Err() io.Writer {
	return std.Err()
}

// Writer returns an io.Writer for injecting our logging into 3rd-party
// libraries
func Writer() *LoggedWriter {
	return std.Writer()
}
