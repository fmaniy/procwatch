package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/procwatch/internal/config"
)

const defaultConfigPath = "procwatch.json"

func main() {
	cfgPath := flag.String("config", defaultConfigPath, "path to JSON config file")
	validate := flag.Bool("validate", false, "validate config and exit")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "procwatch: %v\n", err)
		os.Exit(1)
	}

	if *validate {
		fmt.Printf("Config OK: %d process(es) defined.\n", len(cfg.Processes))
		for _, p := range cfg.Processes {
			hc := "no health-check"
			if p.HealthCheck != nil {
				hc = fmt.Sprintf("health-check @ %s every %v", p.HealthCheck.URL, p.HealthCheck.Interval)
			}
			fmt.Printf("  - %-20s  cmd=%s  restart=%v  %s\n",
				p.Name, p.Command, p.AutoRestart, hc)
		}
		os.Exit(0)
	}

	log.Printf("procwatch starting with log_level=%s, %d process(es)",
		cfg.LogLevel, len(cfg.Processes))

	// Supervisor loop will be wired here in subsequent phases.
	log.Println("procwatch: supervisor not yet implemented — exiting")
}
