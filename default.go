package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger     FieldLogger
	defaultLoggerOnce sync.Once
)

// NewLogger based on the specified log level, defaults to "debug".
// See `New` for more details.
func NewLogger(level string) FieldLogger {
	lvl, err := ParseLevel(level)
	if err != nil {
		lvl = DebugLevel
	}
	return New(lvl)
}

// Default returns the default FieldLogger instance.
// The default logger is initialized with InfoLevel on first call.
func Default() FieldLogger {
	defaultLoggerOnce.Do(func() {
		if defaultLogger == nil {
			defaultLogger = New(InfoLevel)
		}
	})
	return defaultLogger
}

// SetDefault sets the default logger instance.
// Subsequent calls to Default() and package-level logging functions will use this logger.
func SetDefault(l FieldLogger) {
	defaultLogger = l
}

// Debugf logs a formatted debug message using the default logger.
func Debugf(msg string, args ...any) {
	Default().Debugf(msg, args...)
}

// Infof logs a formatted info message using the default logger.
func Infof(msg string, args ...any) {
	Default().Infof(msg, args...)
}

// Printf logs a formatted message using the default logger.
func Printf(msg string, args ...any) {
	Default().Printf(msg, args...)
}

// Warnf logs a formatted warning message using the default logger.
func Warnf(msg string, args ...any) {
	Default().Warnf(msg, args...)
}

// Errorf logs a formatted error message using the default logger.
func Errorf(msg string, args ...any) {
	Default().Errorf(msg, args...)
}

// Fatalf logs a formatted fatal message using the default logger and exits the program.
func Fatalf(msg string, args ...any) {
	Default().Fatalf(msg, args...)
}

// Debug logs a debug message using the default logger.
func Debug(args ...any) {
	Default().Debug(args...)
}

// Info logs an info message using the default logger.
func Info(args ...any) {
	Default().Info(args...)
}

// Warn logs a warning message using the default logger.
func Warn(args ...any) {
	Default().Warn(args...)
}

// Error logs an error message using the default logger.
func Error(args ...any) {
	Default().Error(args...)
}

// Fatal logs a fatal message using the default logger and exits the program.
func Fatal(args ...any) {
	Default().Fatal(args...)
}

// Panic logs a panic message using the default logger and panics.
func Panic(args ...any) {
	Default().Panic(args...)
}

// WithField adds a field to the default logger and returns a new FieldLogger.
func WithField(key string, value any) FieldLogger {
	return Default().WithField(key, value)
}

// WithFields adds multiple fields to the default logger and returns a new FieldLogger.
func WithFields(fields map[string]any) FieldLogger {
	return Default().WithFields(fields)
}

func newDefault(lvl Level) FieldLogger {
	opts := &slog.HandlerOptions{
		Level: toDefaultLevel(lvl),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return &fieldLogger{
		logger: slog.New(handler),
		level:  lvl,
	}
}

// toDefaultLevel converts our Level to slog.Level.
func toDefaultLevel(l Level) slog.Level {
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
