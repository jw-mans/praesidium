package main

import (
	"fmt"
	"os"
	"time"

	"praesidium/pkg/config"
	"praesidium/pkg/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prsdm <command>")
		fmt.Println("Commands: daemon")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "daemon":
		runDaemon()
	default:
		fmt.Printf("Unknown command: %s", os.Args[1])
		os.Exit(1)
	}
}

func runDaemon() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		util.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	util.Info("Starting praesidium daemon . . .")
	util.Info("Using interface: %s", cfg.Iface)
	util.Info("Check interval: %s", cfg.CheckInterval)

	// Just wait indefinitely for demo purposes
	// Monitoring logic would go here

	// It is just a placeholder to keep the daemon running
	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		util.Info("Checking interface %s...", cfg.Iface)
	}
}
