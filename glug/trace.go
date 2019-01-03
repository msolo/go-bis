package glug

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

type Tracer interface {
	StartSpan(name string) Span
}

type Span interface {
	Finish()
}

type span struct {
	name    string
	args    []interface{}
	start   time.Time
	elapsed time.Duration
}

func (ts *span) Finish() {
	ts.elapsed = time.Now().Sub(ts.start)
	fields := map[string]interface{}{
		"duration":    ts.elapsed,
		"durationStr": fmtDuration(ts.elapsed),
		"args":        ts.args,
	}
	msg := fmtLogEntry(ts.name, fields)
	InfofDepth(3, msg, ts.args...)
}

func StartSpan(name string) Span {
	return &span{name: name, start: time.Now()}
}

func Tracef(format string, args ...interface{}) Span {
	return StartSpan(fmt.Sprintf(format, args...))
}

// FIXME(msolo) slow
func fmtLogEntry(format string, p map[string]interface{}) string {
	t := template.Must(template.New("").Parse(format))
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	t.Execute(buf, p)
	return buf.String()
}

func fmtDurationHumanely(d time.Duration) string {
	// Largest time is 2540400h10m10.000000000s
	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	msecs := (u / 1e6) % 1000
	secs := (u / 1e9) % 60
	mins := (u / 60e9) % 60
	hrs := (u / 3600e9) % 24
	days := (u / (24 * 3600e9))
	str := fmt.Sprintf("%dd%dh%dm%d.%03ds", days, hrs, mins, secs, msecs)
	if neg {
		return "-" + str
	}
	return str
}

func fmtDuration(d time.Duration) string {
	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	msecs := (u / 1e6) % 1000
	secs := (u / 1e9)
	str := fmt.Sprintf("%d.%03ds", secs, msecs)
	if neg {
		return "-" + str
	}
	return str
}
