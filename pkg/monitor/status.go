package monitor

import "time"

type Status struct {
	Connected      bool
	Iface          string
	VPNIP          string
	RouteProtected bool
	LastChange     time.Time
}

// TODO: Add ExternalIP and Health
