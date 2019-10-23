package slog

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type lineWriter struct {
	lines []string
}

func (lw lineWriter) LastLine() (line string, ok bool) {
	if len(lw.lines) == 0 {
		return "", false
	}
	return lw.lines[len(lw.lines)-1], false
}

func (lw *lineWriter) Write(b []byte) (int, error) {
	lw.lines = append(lw.lines, string(b))
	return len(b), nil
}

func fakeTime() time.Time {
	// 2015-07-27T11:22Z-05:00
	return time.Unix(1438014120, 0)
}

func init() {
	now = fakeTime
}

func TestStdlogRouting(t *testing.T) {
	cfg := &Config{}
	lw := &lineWriter{}
	h := NewHandler(lw, GlogFmtEntry)
	SetHandler(NewLevelHandler(h, cfg))

	CopyStandardLogTo("WARN")

	msg := "...a plan so cunning you could pin a tail on it and call it a weasel."
	log.Printf(msg)

	lastLine, _ := lw.LastLine()
	if !strings.Contains(lastLine, msg) {
		t.Fatal("stdlog message not present in slog output", lw.lines)
	}
	if !strings.HasPrefix(lastLine, "W") {
		t.Fatalf("stdlog message has wrong level: expected %s, found %c", "W", lastLine[0])
	}
}

func testSlog() (*slogger, *lineWriter) {
	cfg := &Config{}
	lw := &lineWriter{}
	h := NewHandler(lw, GlogFmtEntry)
	h = NewLevelHandler(h, cfg)
	slog := &slogger{h: h, cfg: cfg}
	return slog, lw
}

func TestEntryFields(t *testing.T) {
	slog, lw := testSlog()

	tokenSource := "life.go:42"
	slog.WithSource(tokenSource).Infof("slogger with the source of the answer")
	lastLine, _ := lw.LastLine()
	if !strings.Contains(lastLine, tokenSource) {
		t.Fatalf("tokenSource %s not present in slog output: %s", tokenSource, lastLine)
	}

	tokenSource = "simple error"
	slog.WithError(fmt.Errorf(tokenSource)).Infof("slogger with error")
	lastLine, _ = lw.LastLine()
	if !strings.Contains(lastLine, tokenSource) {
		t.Fatalf("tokenSource %s not present in slog output: %s", tokenSource, lastLine)
	}

	tokenSource = "TestEntryFields"
	slog.WithError(errors.WithStack(fmt.Errorf("simple error with stack"))).Info("slogger with stack error")
	lastLine, _ = lw.LastLine()
	if !strings.Contains(lastLine, tokenSource) {
		t.Fatalf("func name %s not present in slog output, missing stack?: %s", tokenSource, lastLine)
	}
}

func TestFileLine(t *testing.T) {
	slog, lw := testSlog()

	tokenSource := "slog_test.go"

	slog.WithError(nil).Info("check file wrapped slog")
	lastLine, _ := lw.LastLine()
	if !strings.Contains(lastLine, tokenSource) {
		t.Fatalf("file name %s not present in wrapped slog output, missing stack?: %s", tokenSource, lastLine)
	}

	slog.Info("check file direct slog")
	lastLine, _ = lw.LastLine()
	if !strings.Contains(lastLine, tokenSource) {
		t.Fatalf("file name %s not present in direct slog output, missing stack?: %s", tokenSource, lastLine)
	}
}
