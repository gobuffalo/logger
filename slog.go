package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// Slog is a Logger implementation backed by log/slog
type Slog struct {
	logger *slog.Logger
	attrs  []slog.Attr
}

func newSlog(lvl Level) FieldLogger {
	opts := &slog.HandlerOptions{
		Level: toSlogLevel(lvl),
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return &Slog{
		logger: slog.New(handler),
	}
}

// SetOutput sets the output writer for the logger.
func (s *Slog) SetOutput(w io.Writer) {
	opts := &slog.HandlerOptions{}
	s.logger = slog.New(slog.NewTextHandler(w, opts))
}

func (s *Slog) WithField(key string, value any) FieldLogger {
	return s.WithFields(map[string]any{key: value})
}

func (s *Slog) WithFields(m map[string]any) FieldLogger {
	attrs := make([]slog.Attr, 0, len(m)+len(s.attrs))
	attrs = append(attrs, s.attrs...)
	for k, v := range m {
		attrs = append(attrs, slog.Any(k, v))
	}
	return &Slog{
		logger: s.logger,
		attrs:  attrs,
	}
}

func (s *Slog) log(level slog.Level, msg string, args ...any) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	s.logger.Log(context.Background(), level, msg, s.attrsToAny()...)
}

func (s *Slog) attrsToAny() []any {
	result := make([]any, len(s.attrs))
	for i, a := range s.attrs {
		result[i] = a
	}
	return result
}

func (s *Slog) Debugf(msg string, args ...any) {
	s.log(slog.LevelDebug, msg, args...)
}

func (s *Slog) Infof(msg string, args ...any) {
	s.log(slog.LevelInfo, msg, args...)
}

func (s *Slog) Printf(msg string, args ...any) {
	s.Infof(msg, args...)
}

func (s *Slog) Warnf(msg string, args ...any) {
	s.log(slog.LevelWarn, msg, args...)
}

func (s *Slog) Errorf(msg string, args ...any) {
	s.log(slog.LevelError, msg, args...)
}

func (s *Slog) Fatalf(msg string, args ...any) {
	s.log(slog.Level(12), msg, args...)
	os.Exit(1)
}

func (s *Slog) Debug(args ...any) {
	s.log(slog.LevelDebug, fmt.Sprint(args...))
}

func (s *Slog) Info(args ...any) {
	s.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (s *Slog) Warn(args ...any) {
	s.log(slog.LevelWarn, fmt.Sprint(args...))
}

func (s *Slog) Error(args ...any) {
	s.log(slog.LevelError, fmt.Sprint(args...))
}

func (s *Slog) Fatal(args ...any) {
	s.log(slog.Level(12), fmt.Sprint(args...))
	os.Exit(1)
}

func (s *Slog) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	s.log(slog.Level(12), msg)
	panic(msg)
}
