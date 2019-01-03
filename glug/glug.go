package glug

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
	host     = "unknownhost"
	userName = "unknownuser"
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
		host = shortHostname(h)
	}

	current, err := user.Current()
	if err == nil {
		userName = current.Username
	}

	// Sanitize userName since it may contain filepath separators on Windows.
	userName = strings.Replace(userName, `\`, "_", -1)

	CopyStandardLogTo("WARNING")

	// Default stderrThreshold is ERROR.
	logging.stderrThreshold = errorLog

	logging.setVState(0, nil, false)
}

func RegisterFlags(fs *flag.FlagSet) {
	fs.Var(&logging.stderrThreshold, "log.level", "logs at or above this threshold go to stderr")
	fs.Var(&logging.traceLocation, "log.backtrace-at", "when logging hits line file:N, emit a stack trace")
}

func SetLevel(name string) {
	logging.stderrThreshold.Set(name)
}

func InfofDepth(depth int, format string, args ...interface{}) {
	InfoDepth(depth, fmt.Sprintf(format, args...))
}

func WarningfDepth(depth int, format string, args ...interface{}) {
	WarningDepth(depth, fmt.Sprintf(format, args...))
}

func ErrorfDepth(depth int, format string, args ...interface{}) {
	ErrorDepth(depth, fmt.Sprintf(format, args...))
}
