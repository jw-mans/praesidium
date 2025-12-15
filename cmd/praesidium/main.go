package main

import (
	"os"
	"time"

	"praesidium/pkg/config"
	"praesidium/pkg/monitor"
	"praesidium/pkg/util"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "daemon" {
		util.Info("Usage: praesidium daemon")
		return
	}

	runDaemon()
}

func runDaemon() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		util.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	util.Info("Starting praesidium daemon...")
	util.Info("Using interface: %s", cfg.Iface)
	util.Info("Check interval: %s", cfg.CheckInterval)

	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		status, err := monitor.CheckInterface(cfg.Iface)
		if err != nil {
			util.Error("VPN down: %v", err)
		} else {
			extIP, err := monitor.GetExternalIP(cfg.IPCheckURL)
			if err != nil {
				util.Error("Failed to get external IP: %v", err)
			} else {
				status.ExternalIP = extIP
			}

			util.Info("VPN up: %s (%s). IP: %s",
				status.Iface,
				status.VPNIP,
				status.ExternalIP,
			)
		}
	}
}
