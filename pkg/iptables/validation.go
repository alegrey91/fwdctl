package iptables

import (
	"fmt"
	"net"
)

func validateIface(iface string) error {
	if iface == "" {
		return fmt.Errorf("name is empty")
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	found := false
	for _, i := range ifaces {
		if i.Name == iface {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("not found")
	}
	return nil
}

func validateProto(proto string) error {
	if proto == "" {
		return fmt.Errorf("protocol name is empty")
	}
	if (proto != "tcp") && (proto != "udp") && (proto != "icmp") {
		return fmt.Errorf("protocol name not allowed")
	}
	return nil
}

func validatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port number not allowed")
	}
	return nil
}

func validateAddress(address string) error {
	// not a valid check for now.
	if address == "" {
		return fmt.Errorf("address is empty")
	}
	return nil
}

// validate returns both bool and error.
// The boolean return true in case the rule passes all checks.
// In case it does not, then the error will describe the problem.
func validate(iface string, proto string, dport int, saddr string, sport int) error {
	err := validateIface(iface)
	if err != nil {
		return fmt.Errorf("interface: '%s' %v", iface, err)
	}

	err = validateProto(proto)
	if err != nil {
		return fmt.Errorf("protocol: '%s' %v", proto, err)
	}

	err = validatePort(dport)
	if err != nil {
		return fmt.Errorf("destination port: '%d' %v", dport, err)
	}

	err = validateAddress(saddr)
	if err != nil {
		return fmt.Errorf("source address: '%s' %v", saddr, err)
	}

	err = validatePort(sport)
	if err != nil {
		return fmt.Errorf("source port: '%d' %v", sport, err)
	}
	return nil
}
