// Package logger provides a structured logging interface for Go applications.
// It uses the standard library's log/slog package as the default implementation,
// with support for structured fields and multiple log levels.
//
// The package defines Logger and FieldLogger interfaces that are used throughout
// Buffalo apps and other systems. The default implementation writes to stdout
// in a human-readable but parseable format.
//
// Basic usage:
//
//	log := logger.New(logger.InfoLevel)
//	log.Info("Server starting")
//	log.WithField("port", 8080).Info("Server started")
//
// For users who prefer the logrus backend, import the logrus subpackage:
//
//	import "github.com/gobuffalo/logger/logrus"
//	log := logrus.New(logrus.InfoLevel)
package logger

import (
	"fmt"
	"log/slog"
)

type FieldLogger interface {
	Logger
	WithField(string, any) FieldLogger
	WithFields(map[string]any) FieldLogger
}

type Logger interface {
	Debugf(string, ...any)
	Infof(string, ...any)
	Printf(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Fatalf(string, ...any)
	Debug(...any)
	Info(...any)
	Warn(...any)
	Error(...any)
	Fatal(...any)
	Panic(...any)
}

// ParseLevel parses a string level into a Level.
func ParseLevel(level string) (Level, error) {
	switch level {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", level)
}

// NewLogger based on the specified log level, defaults to "debug".
// See `New` for more details.
func NewLogger(level string) FieldLogger {
	lvl, err := ParseLevel(level)
	if err != nil {
		lvl = DebugLevel
	}
	return New(lvl)
}

// New based on the specified log level, defaults to "debug".
// This logger will log to the STDOUT in a human readable,
// but parseable form.
//
//	Example: time="2016-12-01T21:02:07-05:00" level=info duration=225.283µs human_size="106 B" method=GET path="/" render=199.79µs request_id=2265736089 size=106 status=200
func New(lvl Level) FieldLogger {
	return newSlog(lvl)
}

// toSlogLevel converts our Level to slog.Level.
func toSlogLevel(l Level) slog.Level {
	switch l {
	case PanicLevel, FatalLevel:
		return slog.Level(12) // custom level above error
	case ErrorLevel:
		return slog.LevelError
	case WarnLevel:
		return slog.LevelWarn
	case InfoLevel:
		return slog.LevelInfo
	case DebugLevel:
		return slog.LevelDebug
	}
	return slog.LevelDebug
}
