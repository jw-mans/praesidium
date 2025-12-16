package monitor

import (
	"fmt"
	"praesidium/pkg/util"

	"github.com/vishvananda/netlink"
)

func CheckDefaultRoute(ifaceName string) (bool, error) {
	routes, err := netlink.RouteList(nil, util.V4_FAMILY)
	if err != nil {
		return false, fmt.Errorf("failed to list routes: %v", err)
	}

	for _, r := range routes {
		if r.Dst == nil {
			iface, err := netlink.LinkByIndex(r.LinkIndex)
			if err != nil {
				return false, fmt.Errorf("failed to get link for default route: %v", err)
			}
			if iface.Attrs().Name == ifaceName {
				return true, nil
			} else {
				return false, fmt.Errorf("default route is on %s", iface.Attrs().Name)
			}
		}
	}
	return false, fmt.Errorf("no default route found")
}
