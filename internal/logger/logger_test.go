package logger_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/yourorg/procwatch/internal/logger"
)

func decodeEntry(t *testing.T, buf *bytes.Buffer) logger.Entry {
	t.Helper()
	var e logger.Entry
	if err := json.NewDecoder(buf).Decode(&e); err != nil {
		t.Fatalf("failed to decode log entry: %v", err)
	}
	return e
}

func TestInfo(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Info("hello world")

	e := decodeEntry(t, &buf)
	if e.Level != logger.LevelInfo {
		t.Errorf("expected level INFO, got %s", e.Level)
	}
	if e.Message != "hello world" {
		t.Errorf("unexpected message: %s", e.Message)
	}
	if e.Timestamp == "" {
		t.Error("timestamp should not be empty")
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)
	l.Error("something failed", errors.New("boom"))

	e := decodeEntry(t, &buf)
	if e.Level != logger.LevelError {
		t.Errorf("expected level ERROR, got %s", e.Level)
	}
	if !strings.Contains(e.Error, "boom") {
		t.Errorf("expected error field to contain 'boom', got %q", e.Error)
	}
}

func TestWithProcess(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf).WithProcess("myapp")
	l.Info("started")

	e := decodeEntry(t, &buf)
	if e.Process != "myapp" {
		t.Errorf("expected process 'myapp', got %q", e.Process)
	}
}

func TestNilWriterDefaultsToStdout(t *testing.T) {
	// Just ensure New(nil) doesn't panic.
	l := logger.New(nil)
	l.Info("no panic")
}

func TestDebugAndWarn(t *testing.T) {
	var buf bytes.Buffer
	l := logger.New(&buf)

	l.Debug("debug msg")
	e := decodeEntry(t, &buf)
	if e.Level != logger.LevelDebug {
		t.Errorf("expected DEBUG, got %s", e.Level)
	}

	l.Warn("warn msg")
	e = decodeEntry(t, &buf)
	if e.Level != logger.LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
}
