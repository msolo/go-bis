package slog

import (
	"log"
	"os"
)

type logEntry struct {
	Subject string
	Verb    string
	Object  string
}

func (le logEntry) Fields() Fields {
	return map[string]interface{}{
		"Subject": le.Subject,
		"Verb":    le.Verb,
		"Object":  le.Object,
	}
}

func Example() {
	cfg := &Config{}
	h := NewHandler(os.Stderr, GlogFmtEntry)
	SetHandler(NewLevelHandler(h, cfg))
	CopyStandardLogTo("WARN")

	log.Printf("system logger printf")
}
