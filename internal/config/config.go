package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ProcessConfig defines a single supervised process.
type ProcessConfig struct {
	Name        string            `json:"name"`
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env"`
	WorkDir     string            `json:"workdir"`
	AutoRestart bool              `json:"auto_restart"`
	MaxRestarts int               `json:"max_restarts"`
	HealthCheck *HealthCheckConfig `json:"health_check,omitempty"`
}

// HealthCheckConfig defines an HTTP health-check hook.
type HealthCheckConfig struct {
	URL             string        `json:"url"`
	IntervalSeconds int           `json:"interval_seconds"`
	TimeoutSeconds  int           `json:"timeout_seconds"`
	Retries         int           `json:"retries"`
	Interval        time.Duration `json:"-"`
	Timeout         time.Duration `json:"-"`
}

// Config is the top-level procwatch configuration.
type Config struct {
	LogLevel  string          `json:"log_level"`
	LogFile   string          `json:"log_file"`
	Processes []ProcessConfig `json:"processes"`
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	return &cfg, nil
}

func (c *Config) validate() error {
	names := make(map[string]struct{}, len(c.Processes))
	for i, p := range c.Processes {
		if p.Name == "" {
			return fmt.Errorf("config: process[%d] missing name", i)
		}
		if p.Command == "" {
			return fmt.Errorf("config: process %q missing command", p.Name)
		}
		if _, dup := names[p.Name]; dup {
			return fmt.Errorf("config: duplicate process name %q", p.Name)
		}
		names[p.Name] = struct{}{}
	}
	return nil
}

func (c *Config) applyDefaults() {
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	for i := range c.Processes {
		p := &c.Processes[i]
		if p.MaxRestarts == 0 {
			p.MaxRestarts = 5
		}
		if hc := p.HealthCheck; hc != nil {
			if hc.IntervalSeconds == 0 {
				hc.IntervalSeconds = 30
			}
			if hc.TimeoutSeconds == 0 {
				hc.TimeoutSeconds = 5
			}
			if hc.Retries == 0 {
				hc.Retries = 3
			}
			hc.Interval = time.Duration(hc.IntervalSeconds) * time.Second
			hc.Timeout = time.Duration(hc.TimeoutSeconds) * time.Second
		}
	}
}
