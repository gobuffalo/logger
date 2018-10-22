package logger

import "github.com/sirupsen/logrus"

var _ Logger = Logrus{}

// Logrus is a Logger implementation backed by sirupsen/logrus
type Logrus struct {
	logrus.FieldLogger
}

// WithField returns a new Logger with the field added
func (l Logrus) WithField(s string, i interface{}) Logger {
	return Logrus{l.FieldLogger.WithField(s, i)}
}

// WithFields returns a new Logger with the fields added
func (l Logrus) WithFields(m map[string]interface{}) Logger {
	return Logrus{l.FieldLogger.WithFields(m)}
}
