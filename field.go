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
	level  Level
}

// SetOutput sets the output writer for the logger.
func (s *fieldLogger) SetOutput(w io.Writer) {
	opts := &slog.HandlerOptions{
		Level: toDefaultLevel(s.level),
	}
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
		level:  s.level,
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
