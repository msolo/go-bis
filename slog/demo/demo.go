package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/msolo/go-bis/slog"
	"github.com/pkg/errors"
)

type logEntry struct {
	Subject string
	Verb    string
	Object  string
}

func (le logEntry) Fields() slog.Fields {
	return map[string]interface{}{
		"Subject": le.Subject,
		"Verb":    le.Verb,
		"Object":  le.Object,
	}
}

func main() {
	binfile := flag.String("log.binfile", "", "Set up a binary log format.")
	cfg := &slog.Config{}
	slog.RegisterFlags(flag.CommandLine, cfg)
	flag.Parse()

	slog.CopyStandardLogTo("WARN")

	if *binfile != "" {
		f, err := os.OpenFile(*binfile, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		jlogH := slog.NewHandler(f, slog.JsonFmtEntry)
		slog.SetHandler(
			slog.NewLevelHandler(
				slog.NewMultiHandler(slog.GetHandler(), jlogH),
				cfg))
	}

	log.Printf("system logger printf")

	slog.Infof("slogger infof: %s", "OK")
	slog.Warnf("slogger warnf: %s", "OK")
	slog.Errorf("slogger errorf: %s", "OK")

	slog.WithSource("life.go:42").Infof("slogger with the source of the answer")
	slog.WithError(fmt.Errorf("simple error")).Infof("slogger with with error")
	slog.WithError(errors.WithStack(fmt.Errorf("simple error with stack"))).Info("slogger with stack error")

	slog.WithFields(slog.Fields{"field-a": "value-a"}).Infof("slogger with fields")
	slog.WithFields(slog.Fields{"field-a": "values: a\nb\nc"}).Infof("slogger with multiline fields")

	slog.WithFielder(logEntry{"my subject", "is", "objects"}).Info("slogger with fielder")
}
