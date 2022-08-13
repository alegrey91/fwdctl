package iptables

import (
	"net"
)

func interfaceExists(iface string) (bool, error) {
	ifi, err := net.InterfaceByName(iface)
	if err != nil {
		return false, err
	}

	if ifi != nil {
		return true, nil
	}
	return false, nil
}
