package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// fieldLogger is a Logger implementation backed by log/slog
type fieldLogger struct {
	logger *slog.Logger
	attrs  []slog.Attr
}

// SetOutput sets the output writer for the logger.
func (s *fieldLogger) SetOutput(w io.Writer) {
	opts := &slog.HandlerOptions{}
	s.logger = slog.New(slog.NewTextHandler(w, opts))
}

func (s *fieldLogger) WithField(key string, value any) FieldLogger {
	return s.WithFields(map[string]any{key: value})
}

func (s *fieldLogger) WithFields(m map[string]any) FieldLogger {
	attrs := make([]slog.Attr, 0, len(m)+len(s.attrs))
	attrs = append(attrs, s.attrs...)
	for k, v := range m {
		attrs = append(attrs, slog.Any(k, v))
	}
	return &fieldLogger{
		logger: s.logger,
		attrs:  attrs,
	}
}

func (s *fieldLogger) log(level slog.Level, msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	s.logger.Log(context.Background(), level, msg, s.attrsToAny()...)
}

func (s *fieldLogger) attrsToAny() []any {
	result := make([]any, len(s.attrs))
	for i, a := range s.attrs {
		result[i] = a
	}
	return result
}

func (s *fieldLogger) Debugf(msg string, args ...any) {
	s.log(slog.LevelDebug, msg, args...)
}

func (s *fieldLogger) Infof(msg string, args ...any) {
	s.log(slog.LevelInfo, msg, args...)
}

func (s *fieldLogger) Printf(msg string, args ...any) {
	s.Infof(msg, args...)
}

func (s *fieldLogger) Warnf(msg string, args ...any) {
	s.log(slog.LevelWarn, msg, args...)
}

func (s *fieldLogger) Errorf(msg string, args ...any) {
	s.log(slog.LevelError, msg, args...)
}

func (s *fieldLogger) Fatalf(msg string, args ...any) {
	s.log(slog.Level(12), msg, args...)
	os.Exit(1)
}

func (s *fieldLogger) Debug(args ...any) {
	s.log(slog.LevelDebug, fmt.Sprint(args...))
}

func (s *fieldLogger) Info(args ...any) {
	s.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (s *fieldLogger) Warn(args ...any) {
	s.log(slog.LevelWarn, fmt.Sprint(args...))
}

func (s *fieldLogger) Error(args ...any) {
	s.log(slog.LevelError, fmt.Sprint(args...))
}

func (s *fieldLogger) Fatal(args ...any) {
	s.log(slog.Level(12), fmt.Sprint(args...))
	os.Exit(1)
}

func (s *fieldLogger) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	s.log(slog.Level(12), msg)
	panic(msg)
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

func newSlog(lvl Level) FieldLogger {
	opts := &slog.HandlerOptions{
		Level: toSlogLevel(lvl),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return &fieldLogger{
		logger: slog.New(handler),
	}
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
