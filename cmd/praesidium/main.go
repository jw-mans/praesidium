package main

import (
	"os"
	"time"

	"praesidium/pkg/actions"
	"praesidium/pkg/config"
	"praesidium/pkg/killswitch"
	"praesidium/pkg/monitor"
	"praesidium/pkg/server"
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

	util.Info("Starting praesidium daemon")
	util.Info("Interface: %s", cfg.Iface)
	util.Info("Check interval: %s", cfg.CheckInterval)

	// status store
	store := monitor.NewStatusStore(cfg.Iface)

	// HTTP API
	srv := server.New(store)
	srv.Start("127.0.0.1:16888")

	ticker := time.NewTicker(cfg.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		status, err := monitor.CheckInterface(cfg.Iface)
		if err != nil {
			util.Error("VPN down: %v", err)
			status.RouteProtected = false
		} else {
			// external IP
			extIP, err := monitor.GetExternalIP(cfg.IPCheckURL)
			if err != nil {
				util.Error("Failed to get external IP: %v", err)
			} else {
				status.ExternalIP = extIP
			}

			// route check
			protected, err := monitor.CheckDefaultRoute(cfg.Iface)
			if err != nil {
				util.Error("Route leak detected: %v", err)
				status.RouteProtected = false
			} else {
				status.RouteProtected = protected
			}
		}

		// log state
		if status.Connected {
			util.Info(
				"VPN up: %s (%s), external IP: %s, route protected: %v",
				status.Iface,
				status.VPNIP,
				status.ExternalIP,
				status.RouteProtected,
			)
		}

		// update shared status
		store.Update(*status)

		// killswitch + actions
		if !status.Connected || !status.RouteProtected {
			if err := killswitch.Activate(cfg.Iface); err != nil {
				util.Error("Killswitch error: %v", err)
			}

			if len(cfg.OnDisconnect) > 0 {
				actions.RunOnDisconnectActions(cfg.OnDisconnect)
			}
		}
	}
}
