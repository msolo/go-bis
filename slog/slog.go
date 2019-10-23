package slog

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	hostname = "unknownhost"
	username = "unknownuser"
)

// shortHostname returns its argument, truncating at the first period.
// For instance, given "www.google.com" it returns "www".
func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

func init() {
	h, err := os.Hostname()
	if err == nil {
		hostname = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		username = current.Username
	}
}

type Config struct {
	Level Level
	Fname string
}

// Register the flags on the default logger.
func RegisterDefaultFlags() {
	RegisterFlags(flag.CommandLine, std.cfg)
}

func RegisterFlags(fs *flag.FlagSet, cfg *Config) {
	fs.Var(&cfg.Level, "log.level", "logs at or above this threshold")
	fs.StringVar(&cfg.Fname, "log.file", "/dev/stderr", "direct logs to this file")
}

type logHandler struct {
	mu       sync.Mutex
	wr       io.Writer
	fmtEntry FmtEntry
}

func (h *logHandler) WriteEntry(e Entry) error {
	data := h.fmtEntry(e)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.wr, data)
	return err
}

func GlogFmtEntry(e Entry) string {
	dateTime := e.Timestamp().Format("0102 15:04:05")
	micros := e.Timestamp().Nanosecond() / 1e3

	levelName := levelChar[e.Level()]

	fm := Fields{}
	if e.Err() != nil {
		fm["Err"] = e.Err().Error()
	}
	if st := e.StackTrace(); st != nil {
		fm["StackTrace"] = st
	}
	if fields := e.Fields(); len(fields) > 0 {
		fm["Fields"] = fields
	}
	data := []byte("{}")
	var err error
	if len(fm) > 0 {
		data, err = json.MarshalIndent(fm, "", "")
		if err == nil {
			buf := &bytes.Buffer{}
			if err := json.Compact(buf, data); err == nil {
				data = buf.Bytes()
			}
		}
		if err != nil {
			// No point in handling this error.
			data, _ = json.MarshalIndent(map[string]string{"JsonErr": err.Error()}, "", "")
		}
	}

	return fmt.Sprintf("%c%s.%06d %d %s] %s | %s\n",
		levelName, dateTime, micros, pid, e.Source(), e.Message(), data)
}

// Allow override for testing.
var now = time.Now

type slogger struct {
	mu  sync.Mutex
	h   Handler
	cfg *Config
}

func (lg *slogger) Info(args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(InfoLevel, fmt.Sprint(args...))
}

func (lg *slogger) Infof(format string, args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(InfoLevel, fmt.Sprintf(format, args...))
}

func (lg *slogger) Warn(args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(WarnLevel, fmt.Sprint(args...))
}

func (lg *slogger) Warnf(format string, args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(WarnLevel, fmt.Sprintf(format, args...))
}

func (lg *slogger) Error(args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(ErrorLevel, fmt.Sprint(args...))
}

func (lg *slogger) Errorf(format string, args ...interface{}) {
	esl := entrySlogger{handler: lg.h}
	esl.log(ErrorLevel, fmt.Sprintf(format, args...))
}

func (lg *slogger) WithSource(src string) Slogger {
	return &entrySlogger{entry{source: src}, lg.h}
}

func (lg *slogger) WithError(err error) Slogger {
	return &entrySlogger{entry{err: err}, lg.h}
}

func (lg *slogger) WithFields(f Fields) Slogger {
	return &entrySlogger{entry{fields: f}, lg.h}
}

func (lg *slogger) WithFielder(f Fielder) Slogger {
	return &entrySlogger{entry{fielders: []Fielder{f}}, lg.h}
}

func source(depth int) (string, int) {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return file, line
}

func NewHandler(wr io.Writer, fmtEntry FmtEntry) Handler {
	return &logHandler{wr: wr, fmtEntry: fmtEntry}
}

func NewLevelHandler(h Handler, cfg *Config) Handler {
	return &LevelHandler{h, cfg}
}

type LevelHandler struct {
	h   Handler
	cfg *Config
}

func (lh *LevelHandler) WriteEntry(e Entry) error {
	if e.Level() < lh.cfg.Level {
		return nil
	}
	return lh.h.WriteEntry(e)
}

func new(wr io.Writer) *slogger {
	cfg := &Config{}
	return &slogger{
		h:   NewLevelHandler(NewHandler(wr, GlogFmtEntry), cfg),
		cfg: cfg,
	}
}

var (
	std = new(os.Stderr)

	Infof       = std.Infof
	Info        = std.Info
	Warnf       = std.Warnf
	Warn        = std.Warn
	Errorf      = std.Errorf
	Error       = std.Error
	WithFields  = std.WithFields
	WithFielder = std.WithFielder
	WithError   = std.WithError
	WithSource  = std.WithSource
)

func SetHandler(h Handler) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.h = h
}

func GetHandler() Handler {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.h
}
