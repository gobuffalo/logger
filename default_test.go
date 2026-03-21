package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	t.Run("returns non-nil logger", func(t *testing.T) {
		log := Default()
		if log == nil {
			t.Fatal("Default() should not return nil")
		}
	})

	t.Run("returns same instance on multiple calls", func(t *testing.T) {
		// Reset for this test
		d := Default()
		d2 := Default()

		if d != d2 {
			t.Fatal("Default() should return the same instance on multiple calls")
		}
	})
}

func TestSetDefault(t *testing.T) {
	t.Run("changes the default logger", func(t *testing.T) {
		// Save the original default
		original := Default()

		// Create a custom logger with a buffer
		buf := &bytes.Buffer{}
		custom := New(InfoLevel)
		custom.(Outable).SetOutput(buf)

		// Set the custom logger as default
		SetDefault(custom)

		// Verify the default is now our custom logger
		if Default() != custom {
			t.Fatal("SetDefault() did not change the default logger")
		}

		// Test that package-level functions use the custom logger
		Info("test message")
		if !strings.Contains(buf.String(), "test message") {
			t.Fatalf("Expected custom logger to log message, got: %s", buf.String())
		}

		// Restore original
		SetDefault(original)
	})
}

func TestPackageLevelFunctions(t *testing.T) {
	t.Run("Info logs message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Info("info message")
		if !strings.Contains(buf.String(), "info message") {
			t.Fatalf("Expected Info to log message, got: %s", buf.String())
		}
	})

	t.Run("Infof logs formatted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Infof("formatted %s message %d", "test", 42)
		if !strings.Contains(buf.String(), "formatted test message 42") {
			t.Fatalf("Expected Infof to log formatted message, got: %s", buf.String())
		}
	})

	t.Run("Debug does not log when level is Info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Debug("debug message")
		if strings.Contains(buf.String(), "debug message") {
			t.Fatal("Debug should not log when level is Info")
		}
	})

	t.Run("Warn logs message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(WarnLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Warn("warn message")
		if !strings.Contains(buf.String(), "warn message") {
			t.Fatalf("Expected Warn to log message, got: %s", buf.String())
		}
	})

	t.Run("Warnf logs formatted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(WarnLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Warnf("warn %s", "message")
		if !strings.Contains(buf.String(), "warn message") {
			t.Fatalf("Expected Warnf to log formatted message, got: %s", buf.String())
		}
	})

	t.Run("Error logs message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(ErrorLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Error("error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Expected Error to log message, got: %s", buf.String())
		}
	})

	t.Run("Errorf logs formatted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(ErrorLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Errorf("error %s", "message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Expected Errorf to log formatted message, got: %s", buf.String())
		}
	})

	t.Run("Printf logs message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		Printf("print %s", "message")
		if !strings.Contains(buf.String(), "print message") {
			t.Fatalf("Expected Printf to log message, got: %s", buf.String())
		}
	})
}

func TestWithField(t *testing.T) {
	t.Run("adds field to logger", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		WithField("key", "value").Info("message with field")
		if !strings.Contains(buf.String(), "key=value") {
			t.Fatalf("Expected field to be logged, got: %s", buf.String())
		}
	})

	t.Run("returns new FieldLogger", func(t *testing.T) {
		original := Default()
		newLog := WithField("key", "value")

		if original == newLog {
			t.Fatal("WithField should return a new logger instance")
		}
	})
}

func TestWithFields(t *testing.T) {
	t.Run("adds multiple fields to logger", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)
		SetDefault(log)

		WithFields(map[string]any{
			"key1": "value1",
			"key2": 42,
		}).Info("message with fields")

		output := buf.String()
		if !strings.Contains(output, "key1=value1") {
			t.Fatalf("Expected key1 to be logged, got: %s", output)
		}
		if !strings.Contains(output, "key2=42") {
			t.Fatalf("Expected key2 to be logged, got: %s", output)
		}
	})

	t.Run("returns new FieldLogger", func(t *testing.T) {
		original := Default()
		newLog := WithFields(map[string]any{"key": "value"})

		if original == newLog {
			t.Fatal("WithFields should return a new logger instance")
		}
	})
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
		wantErr  bool
	}{
		{"debug", "debug", DebugLevel, false},
		{"info", "info", InfoLevel, false},
		{"warn", "warn", WarnLevel, false},
		{"warning", "warning", WarnLevel, false},
		{"error", "error", ErrorLevel, false},
		{"fatal", "fatal", FatalLevel, false},
		{"panic", "panic", PanicLevel, false},
		{"invalid", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := ParseLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseLevel(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && level != tt.expected {
				t.Fatalf("ParseLevel(%q) = %v, expected %v", tt.input, level, tt.expected)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Run("creates logger with valid level", func(t *testing.T) {
		log := NewLogger("info")
		if log == nil {
			t.Fatal("NewLogger should not return nil")
		}
	})

	t.Run("defaults to debug for invalid level", func(t *testing.T) {
		// This should not panic and should return a logger
		log := NewLogger("invalid")
		if log == nil {
			t.Fatal("NewLogger should not return nil for invalid level")
		}
	})
}

func TestFieldLogger_SetOutput(t *testing.T) {
	t.Run("changes output destination", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.Info("test message")
		if !strings.Contains(buf.String(), "test message") {
			t.Fatalf("Expected message in buffer, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Debug(t *testing.T) {
	t.Run("logs debug message at debug level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(DebugLevel)
		log.(Outable).SetOutput(buf)

		log.Debug("debug message")
		if !strings.Contains(buf.String(), "debug message") {
			t.Fatalf("Expected debug message, got: %s", buf.String())
		}
	})

	t.Run("Debugf logs formatted debug message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(DebugLevel)
		log.(Outable).SetOutput(buf)

		log.Debugf("debug %s %d", "test", 42)
		if !strings.Contains(buf.String(), "debug test 42") {
			t.Fatalf("Expected formatted debug message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Info(t *testing.T) {
	t.Run("logs info message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.Info("info message")
		if !strings.Contains(buf.String(), "info message") {
			t.Fatalf("Expected info message, got: %s", buf.String())
		}
	})

	t.Run("Infof logs formatted message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.Infof("info %s", "formatted")
		if !strings.Contains(buf.String(), "info formatted") {
			t.Fatalf("Expected formatted info message, got: %s", buf.String())
		}
	})

	t.Run("Printf logs as info", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.Printf("print %s", "message")
		if !strings.Contains(buf.String(), "print message") {
			t.Fatalf("Expected print message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Warn(t *testing.T) {
	t.Run("logs warn message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(WarnLevel)
		log.(Outable).SetOutput(buf)

		log.Warn("warn message")
		if !strings.Contains(buf.String(), "warn message") {
			t.Fatalf("Expected warn message, got: %s", buf.String())
		}
	})

	t.Run("Warnf logs formatted warn message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(WarnLevel)
		log.(Outable).SetOutput(buf)

		log.Warnf("warn %s", "formatted")
		if !strings.Contains(buf.String(), "warn formatted") {
			t.Fatalf("Expected formatted warn message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_Error(t *testing.T) {
	t.Run("logs error message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(ErrorLevel)
		log.(Outable).SetOutput(buf)

		log.Error("error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Expected error message, got: %s", buf.String())
		}
	})

	t.Run("Errorf logs formatted error message", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(ErrorLevel)
		log.(Outable).SetOutput(buf)

		log.Errorf("error %s", "formatted")
		if !strings.Contains(buf.String(), "error formatted") {
			t.Fatalf("Expected formatted error message, got: %s", buf.String())
		}
	})
}

func TestFieldLogger_WithField(t *testing.T) {
	t.Run("adds single field", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.WithField("key", "value").Info("message")
		if !strings.Contains(buf.String(), "key=value") {
			t.Fatalf("Expected field in output, got: %s", buf.String())
		}
	})

	t.Run("chains multiple WithField calls", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.WithField("key1", "value1").WithField("key2", "value2").Info("message")
		output := buf.String()
		if !strings.Contains(output, "key1=value1") {
			t.Fatalf("Expected key1 in output, got: %s", output)
		}
		if !strings.Contains(output, "key2=value2") {
			t.Fatalf("Expected key2 in output, got: %s", output)
		}
	})

	t.Run("original logger is unchanged", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		newLog := log.WithField("key", "value")
		log.Info("original message")

		if strings.Contains(buf.String(), "key=value") {
			t.Fatal("Original logger should not have the field")
		}

		newBuf := &bytes.Buffer{}
		newLog.(Outable).SetOutput(newBuf)
		newLog.Info("new message")
		if !strings.Contains(newBuf.String(), "key=value") {
			t.Fatalf("New logger should have the field, got: %s", newBuf.String())
		}
	})
}

func TestFieldLogger_WithFields(t *testing.T) {
	t.Run("adds multiple fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.WithFields(map[string]any{
			"key1": "value1",
			"key2": 42,
		}).Info("message")

		output := buf.String()
		if !strings.Contains(output, "key1=value1") {
			t.Fatalf("Expected key1 in output, got: %s", output)
		}
		if !strings.Contains(output, "key2=42") {
			t.Fatalf("Expected key2 in output, got: %s", output)
		}
	})

	t.Run("combines with existing fields", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.WithField("key1", "value1").WithFields(map[string]any{
			"key2": "value2",
		}).Info("message")

		output := buf.String()
		if !strings.Contains(output, "key1=value1") {
			t.Fatalf("Expected key1 in output, got: %s", output)
		}
		if !strings.Contains(output, "key2=value2") {
			t.Fatalf("Expected key2 in output, got: %s", output)
		}
	})
}

func TestFieldLogger_LevelFiltering(t *testing.T) {
	t.Run("debug does not log at info level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(InfoLevel)
		log.(Outable).SetOutput(buf)

		log.Debug("debug message")
		if strings.Contains(buf.String(), "debug message") {
			t.Fatal("Debug should not log at InfoLevel")
		}
	})

	t.Run("info does not log at warn level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(WarnLevel)
		log.(Outable).SetOutput(buf)

		log.Info("info message")
		if strings.Contains(buf.String(), "info message") {
			t.Fatal("Info should not log at WarnLevel")
		}
	})

	t.Run("warn does not log at error level", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(ErrorLevel)
		log.(Outable).SetOutput(buf)

		log.Warn("warn message")
		if strings.Contains(buf.String(), "warn message") {
			t.Fatal("Warn should not log at ErrorLevel")
		}
	})

	t.Run("higher levels log at lower levels", func(t *testing.T) {
		buf := &bytes.Buffer{}
		log := New(DebugLevel)
		log.(Outable).SetOutput(buf)

		log.Error("error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Fatalf("Error should log at DebugLevel, got: %s", buf.String())
		}
	})
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warning"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{PanicLevel, "panic"},
		{Level(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Fatalf("Level(%d).String() = %q, expected %q", tt.level, result, tt.expected)
			}
		})
	}
}
