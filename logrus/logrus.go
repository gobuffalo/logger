// Package logrus provides a logger implementation backed by sirupsen/logrus.
//
// This is an optional subpackage for users who prefer logrus over the default
// slog-based implementation in the root logger package. Importing this package
// adds logrus as a dependency to your project.
//
// Features:
//   - Full compatibility with the logger.FieldLogger interface
//   - Structured logging with fields via WithField/WithFields
//   - Multiple log levels (Debug, Info, Warn, Error, Fatal, Panic)
//   - Formatted logging methods (Debugf, Infof, etc.)
//   - Custom text formatter with color support for terminal output
//
// When to use logrus vs default:
//   - Use logrus if you need hooks, formatters, or logrus-specific features
//   - Use the default (slog-based) for zero external dependencies
//
// Level mapping:
//   - logger.DebugLevel -> logrus.DebugLevel
//   - logger.InfoLevel  -> logrus.InfoLevel
//   - logger.WarnLevel  -> logrus.WarnLevel
//   - logger.ErrorLevel -> logrus.ErrorLevel
//   - logger.FatalLevel -> logrus.FatalLevel
//   - logger.PanicLevel -> logrus.PanicLevel
//
// Usage:
//
//	import "github.com/gobuffalo/logger/v2/logrus"
//
//	// Create a new logger
//	log := logrus.New(logrus.InfoLevel)
//	log.Info("Application started")
//
//	// With structured fields
//	log.WithField("version", "1.0.0").Info("Version info")
//	log.WithFields(map[string]any{
//	    "request_id": "abc123",
//	    "user_id": 42,
//	}).Info("Request processed")
//
//	// From string level (defaults to debug on invalid)
//	log := logrus.NewLogger("info")
//
// The fieldLogger type implements the logger.FieldLogger interface, making it
// fully interchangeable with the default logger implementation. You can pass
// logrus loggers anywhere that accepts logger.FieldLogger.
package logrus

import (
	"io"

	"github.com/gobuffalo/logger/v2"
	"github.com/sirupsen/logrus"
)

var _ logger.Logger = fieldLogger{}
var _ logger.FieldLogger = fieldLogger{}
var _ logger.Outable = fieldLogger{}

// fieldLogger is a Logger implementation backed by sirupsen/logrus
type fieldLogger struct {
	logrus.FieldLogger
}

// SetOutput will try and set the output of the underlying
// logrus.FieldLogger if it can
func (l fieldLogger) SetOutput(w io.Writer) {
	if lg, ok := l.FieldLogger.(interface{ SetOutput(io.Writer) }); ok {
		lg.SetOutput(w)
	}
}

func (l fieldLogger) WithField(s string, i any) logger.FieldLogger {
	return fieldLogger{l.FieldLogger.WithField(s, i)}
}

func (l fieldLogger) WithFields(m map[string]any) logger.FieldLogger {
	return fieldLogger{l.FieldLogger.WithFields(logrus.Fields(m))}
}

// New creates a new logger.FieldLogger backed by logrus.
func New(lvl logger.Level) logger.FieldLogger {
	l := logrus.New()
	l.SetLevel(logrus.Level(lvl))
	return fieldLogger{l}
}

// NewLogger creates a new logger.FieldLogger from a string level.
func NewLogger(level string) logger.FieldLogger {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.DebugLevel
	}
	return New(logger.Level(lvl))
}
