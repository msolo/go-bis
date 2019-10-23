package slog

import (
	"bytes"
	"fmt"
	stdLog "log"
	"strconv"
)

// CopyStandardLogTo arranges for messages written to the Go "log" package's
// default logs to also appear in the strucuted logs for the named and lower
// severities.
//
// Valid names are "INFO", "WARNING", "ERROR".  If the name is not
// recognized, CopyStandardLogTo panics.
func CopyStandardLogTo(name string) {
	sev, err := parseLevel(name)
	if err != nil {
		panic(err)
	}
	// Set a log format that captures the user's file and line:
	//   d.go:23: message
	stdLog.SetFlags(stdLog.Lshortfile)
	stdLog.SetOutput(logBridge(sev))
}

// logBridge provides the Write method that enables CopyStandardLogTo to connect
// Go's standard logs to the logs provided by this package.
type logBridge Level

// Write parses the standard logging line and passes its components to the
// logger for severity(lb).
// FIXME(msolo) log messages with newlines should be encoded to prevent hijacking the overall
// log format. I suppose the real solution is to use JSON or binary logs.
func (lb logBridge) Write(b []byte) (n int, err error) {
	var (
		file = "???"
		line = 1
		text string
	)
	b = bytes.TrimSpace(b)
	// Split "d.go:23: message" into "d.go", "23", and "message".
	if parts := bytes.SplitN(b, []byte{':'}, 3); len(parts) != 3 || len(parts[0]) < 1 || len(parts[2]) < 1 {
		text = fmt.Sprintf("bad log format: %s", b)
	} else {
		file = string(parts[0])
		text = string(parts[2][1:]) // skip leading space
		line, err = strconv.Atoi(string(parts[1]))
		if err != nil {
			text = fmt.Sprintf("bad line number: %s", b)
			line = 1
		}
	}
	src := fmt.Sprintf("%s:%d", file, line)
	std.WithSource(src).(*entrySlogger).log(Level(lb), string(text))
	return len(text), nil
}
