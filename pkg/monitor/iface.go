package monitor

import (
	"fmt"
	"time"

	"praesidium/pkg/util"

	"github.com/vishvananda/netlink"
)

func CheckInterface(ifaceName string) (*Status, error) {
	status := &Status{
		Iface:      ifaceName,
		LastChange: time.Now(),
	}

	link, err := netlink.LinkByName(ifaceName)
	if err != nil {
		status.Connected = false
		return status, fmt.Errorf("interface not found: %v", err)
	}

	status.Connected = true

	addrs, err := netlink.AddrList(link, util.V4_FAMILY)
	if err != nil || len(addrs) == 0 {
		status.VPNIP = ""
	} else {
		status.VPNIP = addrs[0].IP.String()
	}

	return status, nil
}
