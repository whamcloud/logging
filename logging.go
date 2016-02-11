package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// LogFileFlags are suitable for appending to a log file
	LogFileFlags = os.O_CREATE | os.O_APPEND | os.O_RDWR

	// LogFileMode is suitable for root-only log access
	LogFileMode = 0600
)

// CreateWriter is a convenience function to ensure that the given input
// results in an io.Writer
func CreateWriter(w interface{}) (io.Writer, error) {
	switch w := w.(type) {
	case io.Writer:
		return w, nil
	case string:
		switch strings.ToLower(w) {
		case "stderr":
			return os.Stderr, nil
		case "stdout":
			return os.Stdout, nil
		case "":
			return ioutil.Discard, nil
		default:
			return os.OpenFile(w, LogFileFlags, LogFileMode)
		}
	default:
		return nil, fmt.Errorf("CreateWriter() called with unhandled input: %v", w)
	}
}
