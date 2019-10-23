package slog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type entry struct {
	timeStarted time.Time
	timeEnded   time.Time
	level       Level
	source      string
	message     string
	fields      Fields
	fielders    []Fielder
	err         error
	pid         int
	hostname    string
}

func (ent *entry) MarshalJSON() ([]byte, error) {
	st := struct {
		Level      Level
		Timestamp  time.Time
		Hostname   string
		Pid        int
		Source     string
		Message    string
		Fields     Fields            `json:",omitempty"`
		Err        string            `json:",omitempty"`
		StackTrace errors.StackTrace `json:",omitempty"`
	}{ent.level,
		ent.timeStarted,
		ent.hostname,
		ent.pid,
		ent.source,
		ent.message,
		ent.Fields(),
		maybeErrString(ent.err),
		ent.StackTrace(),
	}
	return json.Marshal(st)
}

func (ent *entry) Timestamp() time.Time {
	return ent.timeStarted
}

func (ent *entry) Err() error {
	return ent.err
}

func (ent *entry) Source() string {
	return ent.source
}

func (ent *entry) Level() Level {
	return ent.level
}

func (ent *entry) Message() string {
	return ent.message
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (ent *entry) StackTrace() errors.StackTrace {
	if err, ok := ent.err.(stackTracer); ok {
		return err.StackTrace()
	}
	return nil
}

func maybeErrString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func (ent *entry) Fields() Fields {
	var mf Fields
	for _, fielder := range ent.fielders {
		mf = mergeFields(mf, fielder.Fields())
	}
	return mergeFields(mf, ent.fields)
}

func (ent *entry) Pid() int {
	return ent.pid
}

func (ent *entry) Hostname() string {
	return ent.hostname
}

type entrySlogger struct {
	entry
	handler Handler
}

func (esl *entrySlogger) WithSource(src string) Slogger {
	esl2 := entrySlogger{}
	esl2 = *esl
	esl2.source = src
	return &esl2
}

func (esl *entrySlogger) WithError(err error) Slogger {
	esl2 := entrySlogger{}
	esl2 = *esl
	esl2.err = err
	return &esl2
}

func (esl *entrySlogger) WithFielder(f Fielder) Slogger {
	esl2 := entrySlogger{}
	esl2 = *esl
	fs := make([]Fielder, 0, 4)
	copy(fs, esl.fielders)
	fs = append(fs, f)
	esl2.fielders = fs
	return &esl2
}

func mergeFields(f1, f2 Fields) Fields {
	mf := make(Fields, len(f1)+len(f2))
	for k, v := range f1 {
		mf[k] = v
	}
	for k, v := range f2 {
		mf[k] = v
	}
	return mf
}

func (esl *entrySlogger) WithFields(f Fields) Slogger {
	esl2 := entrySlogger{}
	esl2 = *esl
	esl2.fields = mergeFields(esl.Fields(), f)
	return &esl2
}

func (esl *entrySlogger) Info(args ...interface{}) {
	esl.log(InfoLevel, fmt.Sprint(args...))
}

func (esl *entrySlogger) Infof(format string, args ...interface{}) {
	esl.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (esl *entrySlogger) Warn(args ...interface{}) {
	esl.log(WarnLevel, fmt.Sprint(args...))
}

func (esl *entrySlogger) Warnf(format string, args ...interface{}) {
	esl.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (esl *entrySlogger) Error(args ...interface{}) {
	esl.log(ErrorLevel, fmt.Sprint(args...))
}

func (esl *entrySlogger) Errorf(format string, args ...interface{}) {
	esl.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (esl *entrySlogger) log(level Level, msg string) {
	ent := &esl.entry
	ent.timeStarted = now().UTC()
	ent.level = level
	ent.message = msg
	ent.pid = pid
	ent.hostname = hostname

	if ent.source == "" {
		file, line := source(0)
		ent.source = fmt.Sprintf("%s:%d", file, line)
	}
	if err := esl.handler.WriteEntry(ent); err != nil {
		println("log write failed:", err)
	}
}
