package logrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gobuffalo/logger/v2"
)

func newTestLogger(t *testing.T, lvl logger.Level, buf *bytes.Buffer) logger.FieldLogger {
	t.Helper()
	log := New(lvl)
	log.(logger.Outable).SetOutput(buf)
	return log
}

func TestNew(t *testing.T) {
	t.Run("creates logger with valid level", func(t *testing.T) {
		log := New(logger.InfoLevel)
		if log == nil {
			t.Fatal("New should not return nil")
		}
	})

	t.Run("implements FieldLogger interface", func(t *testing.T) {
		log := New(logger.DebugLevel)
		var _ logger.FieldLogger = log
	})

	t.Run("implements Outable interface", func(t *testing.T) {
		log := New(logger.DebugLevel)
		var _ logger.Outable = log.(fieldLogger)
	})
}

func TestNewLogger(t *testing.T) {
	t.Run("creates logger from string level", func(t *testing.T) {
		log := NewLogger("info")
		if log == nil {
			t.Fatal("NewLogger should not return nil")
		}
	})

	t.Run("defaults to debug for invalid level", func(t *testing.T) {
		log := NewLogger("invalid")
		if log == nil {
			t.Fatal("NewLogger should not return nil for invalid level")
		}
	})

	t.Run("parses all valid levels", func(t *testing.T) {
		levels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
		for _, lvl := range levels {
			log := NewLogger(lvl)
			if log == nil {
				t.Fatalf("NewLogger(%q) should not return nil", lvl)
			}
		}
	})
}

func TestFieldLogger_SetOutput(t *testing.T) {
	t.Run("changes output destination", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Info("test message")
		if !strings.Contains(buf.String(), "test message") {
			t.Fatalf("Expected message in buffer, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_WithField(t *testing.T) {
	t.Run("adds single field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.WithField("key", "value").Info("message")
		if !strings.Contains(buf.String(), "key") {
			t.Fatalf("Expected field in output, got: %s", buf.String())
		}
	})

	t.Run("returns new FieldLogger", func(t *testing.T) {
		log := New(logger.InfoLevel)
		newLog := log.WithField("key", "value")

		if log == newLog {
			t.Fatal("WithField should return a new logger instance")
		}
	})

	t.Run("chains multiple WithField calls", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.WithField("key1", "value1").WithField("key2", "value2").Info("message")
		output := buf.String()
		if !strings.Contains(output, "key1") || !strings.Contains(output, "key2") {
			t.Fatalf("Expected both fields in output, got: %s", output)
		}
	})
}

func TestFieldLogger_WithFields(t *testing.T) {
	t.Run("adds multiple fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.WithFields(map[string]any{
			"key1": "value1",
			"key2": 42,
		}).Info("message")

		output := buf.String()
		if !strings.Contains(output, "key1") {
			t.Fatalf("Expected key1 in output, got: %s", output)
		}
		if !strings.Contains(output, "key2") {
			t.Fatalf("Expected key2 in output, got: %s", output)
		}
	})

	t.Run("returns new FieldLogger", func(t *testing.T) {
		log := New(logger.InfoLevel)
		newLog := log.WithFields(map[string]any{"key": "value"})

		if log == newLog {
			t.Fatal("WithFields should return a new logger instance")
		}
	})
}

func TestFieldLogger_Debug(t *testing.T) {
	t.Run("logs debug message at debug level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.DebugLevel, buf)

		log.Debug("debug message")
		if !strings.Contains(buf.String(), "debug message") {
			t.Fatalf("Expected debug message, got: %s", buf.String())
		}
	})

	t.Run("Debug does not log at info level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Debug("debug message")
		if strings.Contains(buf.String(), "debug message") {
			t.Fatal("Debug should not log at InfoLevel")
		}
	})

	t.Run("Debugf logs formatted debug message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.DebugLevel, buf)

		log.Debugf("debug %s %d", "test", 42)
		if !strings.Contains(buf.String(), "debug test 42") {
			t.Fatalf("Expected formatted debug message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Info(t *testing.T) {
	t.Run("logs info message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Info("info message")
		if !strings.Contains(buf.String(), "info message") {
			t.Fatalf("Expected info message, got: %s", buf.String())
		}
	})

	t.Run("Infof logs formatted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Infof("info %s", "formatted")
		if !strings.Contains(buf.String(), "info formatted") {
			t.Fatalf("Expected formatted info message, got: %s", buf.String())
		}
	})

	t.Run("Printf logs as info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Printf("print %s", "message")
		if !strings.Contains(buf.String(), "print message") {
			t.Fatalf("Expected print message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Warn(t *testing.T) {
	t.Run("logs warn message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.WarnLevel, buf)

		log.Warn("warn message")
		if !strings.Contains(buf.String(), "warn message") {
			t.Fatalf("Expected warn message, got: %s", buf.String())
		}
	})

	t.Run("Warnf logs formatted warn message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.WarnLevel, buf)

		log.Warnf("warn %s", "formatted")
		if !strings.Contains(buf.String(), "warn formatted") {
			t.Fatalf("Expected formatted warn message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Error(t *testing.T) {
	t.Run("logs error message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.ErrorLevel, buf)

		log.Error("error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Expected error message, got: %s", buf.String())
		}
	})

	t.Run("Errorf logs formatted error message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.ErrorLevel, buf)

		log.Errorf("error %s", "formatted")
		if !strings.Contains(buf.String(), "error formatted") {
			t.Fatalf("Expected formatted error message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_LevelFiltering(t *testing.T) {
	t.Run("debug does not log at info level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.InfoLevel, buf)

		log.Debug("debug message")
		if strings.Contains(buf.String(), "debug message") {
			t.Fatal("Debug should not log at InfoLevel")
		}
	})

	t.Run("info does not log at warn level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.WarnLevel, buf)

		log.Info("info message")
		if strings.Contains(buf.String(), "info message") {
			t.Fatal("Info should not log at WarnLevel")
		}
	})

	t.Run("warn does not log at error level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.ErrorLevel, buf)

		log.Warn("warn message")
		if strings.Contains(buf.String(), "warn message") {
			t.Fatal("Warn should not log at ErrorLevel")
		}
	})

	t.Run("higher levels log at lower levels", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := newTestLogger(t, logger.DebugLevel, buf)

		log.Error("error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Error should log at DebugLevel, got: %s", buf.String())
		}
	})
}
