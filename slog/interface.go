package slog

import (
	"time"

	"github.com/pkg/errors"
)

type Level int

type Fields map[string]interface{}

type Fielder interface {
	Fields() Fields
}

type Entry interface {
	Timestamp() time.Time
	Source() string
	Message() string
	Fields() Fields
	Err() error
	StackTrace() errors.StackTrace
	Pid() int         // too specific?
	Hostname() string // too specific?
	Level() Level
}

type Tracer interface {
	Stop(err error)
}

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	// Fatal(args ...interface{})
	// Fatalf(format string, args ...interface{})

	// Trace(args ...interface{}) Tracer
	// Tracef(format string, args ...interface{}) Tracer
}

type Slogger interface {
	WithFields(f Fields) Slogger
	WithFielder(f Fielder) Slogger
	WithError(err error) Slogger
	WithSource(src string) Slogger
	Logger
}

// Format an entry into a simple string.
type FmtEntry func(e Entry) string

// Write an entry to storage or a stream.
type Handler interface {
	WriteEntry(e Entry) error
}
