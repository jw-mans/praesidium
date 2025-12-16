package killswitch

import (
	"fmt"
	"praesidium/pkg/util"

	"github.com/vishvananda/netlink"
)

func Activate(vpnIface string) error {
	routes, err := netlink.RouteList(nil, util.V4_FAMILY)
	if err != nil {
		return fmt.Errorf("killswitch: failed to list routes: %w", err)
	}
	for _, r := range routes {
		if r.Dst == nil {
			iface, err := netlink.LinkByIndex(r.LinkIndex)
			if err != nil {
				return fmt.Errorf("killswitch: failed to get link for default route: %w", err)
			}

			// already protected
			if iface.Attrs().Name == vpnIface {
				return nil
			}

			// delete leaked default route
			if err := netlink.RouteDel(&r); err != nil {
				return fmt.Errorf("killswitch: failed to delete default route: %w", err)
			}

			util.Error(
				"Killswitch activated: default route removed (was on %s)",
				iface.Attrs().Name,
			)
			return nil
		}
	}

	return fmt.Errorf("killswitch: no default route found")
}
