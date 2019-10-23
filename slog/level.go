package slog

import (
	"fmt"
	"strings"
)

const (
	InvalidLevel Level = iota - 1
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	MaxLevels
)

var levelChar = [MaxLevels]byte{
	'D',
	'I',
	'W',
	'E',
	'F',
}

var levelName = [MaxLevels]string{
	"debug",
	"info",
	"warn",
	"error",
	"fatal",
}

var levelNameMap = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

func parseLevel(val string) (Level, error) {
	x, ok := levelNameMap[strings.ToLower(val)]
	if !ok {
		return InvalidLevel, fmt.Errorf("invalid log level: %s", val)
	}
	return x, nil
}

func (l *Level) Set(val string) (err error) {
	*l, err = parseLevel(val)
	return err
}

func (l *Level) String() string {
	return levelName[int(*l)]
}
