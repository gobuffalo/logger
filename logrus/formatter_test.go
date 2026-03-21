package logrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestTextFormatter_Format(t *testing.T) {
	t.Run("formats basic entry with level and message", func(t *testing.T) {
		f := &textFormatter{}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "test message",
			Level:   logrus.InfoLevel,
			Data:    logrus.Fields{},
		}

		result, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		output := string(result)
		if !strings.Contains(output, "test message") {
			t.Fatalf("Expected message in output, got: %s", output)
		}
		if !strings.Contains(output, "info") && !strings.Contains(output, "INFO") {
			t.Fatalf("Expected level in output, got: %s", output)
		}
	})

	t.Run("includes fields in output", func(t *testing.T) {
		f := &textFormatter{}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "test message",
			Level:   logrus.InfoLevel,
			Data: logrus.Fields{
				"key1": "value1",
				"key2": "value2",
			},
		}

		result, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		output := string(result)
		if !strings.Contains(output, "key1") {
			t.Fatalf("Expected key1 in output, got: %s", output)
		}
		if !strings.Contains(output, "key2") {
			t.Fatalf("Expected key2 in output, got: %s", output)
		}
	})

	t.Run("applies colors when ForceColors is true", func(t *testing.T) {
		f := &textFormatter{ForceColors: true}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "colored message",
			Level:   logrus.ErrorLevel,
			Data:    logrus.Fields{},
		}

		result, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		output := string(result)
		if !strings.Contains(output, "\x1b[") {
			t.Fatalf("Expected ANSI color codes in output, got: %s", output)
		}
	})

	t.Run("quotes values containing spaces", func(t *testing.T) {
		f := &textFormatter{}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "test",
			Level:   logrus.InfoLevel,
			Data: logrus.Fields{
				"key": "value with spaces",
			},
		}

		result, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		output := string(result)
		if !strings.Contains(output, `"value with spaces"`) {
			t.Fatalf("Expected quoted value in output, got: %s", output)
		}
	})

	t.Run("does not quote simple values", func(t *testing.T) {
		f := &textFormatter{}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "test",
			Level:   logrus.InfoLevel,
			Data: logrus.Fields{
				"key": "simple_value-123",
			},
		}

		result, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		output := string(result)
		if strings.Contains(output, `"simple_value-123"`) {
			t.Fatalf("Expected unquoted value in output, got: %s", output)
		}
	})

	t.Run("uses entry buffer when available", func(t *testing.T) {
		f := &textFormatter{}
		buf := &bytes.Buffer{}
		entry := &logrus.Entry{
			Logger:  logrus.New(),
			Message: "buffered",
			Level:   logrus.InfoLevel,
			Data:    logrus.Fields{},
			Buffer:  buf,
		}

		_, err := f.Format(entry)
		if err != nil {
			t.Fatalf("Format returned error: %v", err)
		}

		if buf.Len() == 0 {
			t.Fatal("Expected buffer to contain data")
		}
	})
}

func TestPrefixFieldClashes(t *testing.T) {
	t.Run("renames reserved field 'time'", func(t *testing.T) {
		data := logrus.Fields{"time": "2023-01-01"}
		prefixFieldClashes(data)

		if _, ok := data["fields.time"]; !ok {
			t.Fatal("Expected 'time' to be renamed to 'fields.time'")
		}
	})

	t.Run("renames reserved field 'msg'", func(t *testing.T) {
		data := logrus.Fields{"msg": "hello"}
		prefixFieldClashes(data)

		if _, ok := data["fields.msg"]; !ok {
			t.Fatal("Expected 'msg' to be renamed to 'fields.msg'")
		}
	})

	t.Run("renames reserved field 'level'", func(t *testing.T) {
		data := logrus.Fields{"level": "debug"}
		prefixFieldClashes(data)

		if _, ok := data["fields.level"]; !ok {
			t.Fatal("Expected 'level' to be renamed to 'fields.level'")
		}
	})
}
