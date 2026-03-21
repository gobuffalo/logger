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
//	import "github.com/gobuffalo/logger/v2/logrus"
//	log := logrus.New(logrus.InfoLevel)
package logger

import (
	"io"
)

// Outable interface for loggers that allow setting the output writer
type Outable interface {
	SetOutput(out io.Writer)
}

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
