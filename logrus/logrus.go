// Package logrus provides a logger implementation backed by sirupsen/logrus.
//
// This is an optional subpackage for users who prefer logrus over the default
// slog-based implementation in the root logger package. Importing this package
// adds logrus as a dependency to your project.
//
// Usage:
//
//	import "github.com/gobuffalo/logger/logrus"
//	log := logrus.New(logrus.InfoLevel)
//	log.Info("Application started")
//	log.WithField("version", "1.0.0").Info("Version info")
//
// The Logrus type implements the logger.FieldLogger interface, making it
// interchangeable with the default logger implementation.
package logrus

import (
	"io"

	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

var _ logger.Logger = Logrus{}
var _ logger.FieldLogger = Logrus{}
var _ logger.Outable = Logrus{}

// Logrus is a Logger implementation backed by sirupsen/logrus
type Logrus struct {
	logrus.FieldLogger
}

// SetOutput will try and set the output of the underlying
// logrus.FieldLogger if it can
func (l Logrus) SetOutput(w io.Writer) {
	if lg, ok := l.FieldLogger.(interface{ SetOutput(io.Writer) }); ok {
		lg.SetOutput(w)
	}
}

func (l Logrus) WithField(s string, i any) logger.FieldLogger {
	return Logrus{l.FieldLogger.WithField(s, i)}
}

func (l Logrus) WithFields(m map[string]any) logger.FieldLogger {
	return Logrus{l.FieldLogger.WithFields(logrus.Fields(m))}
}

// New creates a new Logrus logger.
func New(lvl logger.Level) logger.FieldLogger {
	l := logrus.New()
	l.SetLevel(logrus.Level(lvl))
	return Logrus{l}
}

// NewLogger creates a new Logrus logger from a string level.
func NewLogger(level string) logger.FieldLogger {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.DebugLevel
	}
	return New(logger.Level(lvl))
}
