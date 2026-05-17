// Package logger provides a lightweight structured JSON logger for procwatch.
//
// Each log entry is written as a single JSON object on its own line, making
// output easy to consume with tools like jq or any log-aggregation pipeline.
//
// Basic usage:
//
//	l := logger.New(os.Stdout)
//	l.Info("supervisor started")
//
// Scoping a logger to a managed process:
//
//	pl := l.WithProcess("worker")
//	pl.Info("process started")
//	pl.Error("process exited unexpectedly", err)
//
// Log levels available: DEBUG, INFO, WARN, ERROR.
package logger
