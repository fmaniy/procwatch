package logger

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

// Level represents a log severity level.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

// Entry is a structured JSON log entry.
type Entry struct {
	Timestamp string `json:"timestamp"`
	Level     Level  `json:"level"`
	Process   string `json:"process,omitempty"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
}

// Logger writes structured JSON log entries to an io.Writer.
type Logger struct {
	writer  io.Writer
	process string
}

// New creates a Logger that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{writer: w}
}

// WithProcess returns a new Logger scoped to a named process.
func (l *Logger) WithProcess(name string) *Logger {
	return &Logger{writer: l.writer, process: name}
}

func (l *Logger) write(level Level, msg string, err error) {
	e := Entry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Process:   l.process,
		Message:   msg,
	}
	if err != nil {
		e.Error = err.Error()
	}
	data, _ := json.Marshal(e)
	data = append(data, '\n')
	_, _ = l.writer.Write(data)
}

func (l *Logger) Debug(msg string) { l.write(LevelDebug, msg, nil) }
func (l *Logger) Info(msg string)  { l.write(LevelInfo, msg, nil) }
func (l *Logger) Warn(msg string)  { l.write(LevelWarn, msg, nil) }

func (l *Logger) Error(msg string, err error) { l.write(LevelError, msg, err) }
