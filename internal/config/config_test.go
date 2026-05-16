package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "procwatch.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `{
		"log_level": "debug",
		"processes": [
			{
				"name": "web",
				"command": "/usr/bin/python3",
				"args": ["-m", "http.server"],
				"auto_restart": true,
				"health_check": {
					"url": "http://localhost:8000/health",
					"interval_seconds": 10,
					"timeout_seconds": 3,
					"retries": 2
				}
			}
		]
	}`
	cfg, err := Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("log_level: got %q, want %q", cfg.LogLevel, "debug")
	}
	if len(cfg.Processes) != 1 {
		t.Fatalf("expected 1 process, got %d", len(cfg.Processes))
	}
	p := cfg.Processes[0]
	if p.Name != "web" {
		t.Errorf("name: got %q", p.Name)
	}
	if p.HealthCheck.Interval != 10*time.Second {
		t.Errorf("interval: got %v", p.HealthCheck.Interval)
	}
}

func TestLoad_Defaults(t *testing.T) {
	raw := `{"processes":[{"name":"svc","command":"/bin/sh"}]}`
	cfg, err := Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("default log_level: got %q", cfg.LogLevel)
	}
	if cfg.Processes[0].MaxRestarts != 5 {
		t.Errorf("default max_restarts: got %d", cfg.Processes[0].MaxRestarts)
	}
}

func TestLoad_MissingName(t *testing.T) {
	raw := `{"processes":[{"command":"/bin/sh"}]}`
	_, err := Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestLoad_DuplicateName(t *testing.T) {
	raw := `{"processes":[
		{"name":"svc","command":"/bin/sh"},
		{"name":"svc","command":"/bin/bash"}
	]}`
	_, err := Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/procwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
