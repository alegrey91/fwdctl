package iptables

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/alegrey91/fwdctl/internal/extractor"
)

var (
	label string = "fwdctl"
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

func Validate(iface string, proto string, dport int, saddr string, sport int) error {
	return validateForward(iface, proto, dport, saddr, sport)
}

// validateForward returns both bool and error.
// The boolean return true in case the rule passes all checks.
// In case it does not, then the error will describe the problem.
func validateForward(iface string, proto string, dport int, saddr string, sport int) error {
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

func CreateForward(iface string, proto string, dport int, saddr string, sport int) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	// example rule:
	// iptables -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	ruleSpec := []string{
		"-i", iface,
		"-p", proto,
		"-m", proto,
		"--dport", strconv.Itoa(dport),
		"-j", FwdTarget,
		"--to-destination", saddr + ":" + strconv.Itoa(sport),
		"-m", "comment",
		"--comment", label,
	}

	// check if input interface exists on the system
	ifaceExits, err := interfaceExists(iface)
	if err != nil {
		return fmt.Errorf("error reading interfaces: %v", err)
	}
	if !ifaceExits {
		return fmt.Errorf("interface %s does not exists", iface)
	}

	// check if provided rule already exists
	ruleExists, err := ipt.Exists(FwdTable, FwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if ruleExists {
		return fmt.Errorf("rule already exists")
	}

	// apply provided rule
	err = ipt.AppendUnique(FwdTable, FwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("rule failed: %v", err)
	}
	return nil
}

func ListForward(outputFormat string) (map[int]string, error) {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return nil, fmt.Errorf("failed: %v", err)
	}

	ruleList, err := ipt.List(FwdTable, FwdChain)
	if err != nil {
		return nil, fmt.Errorf("failed listing rules: %v", err)
	}

	// check listed rules are tagged with custom tag
	fwdRules := make(map[int]string)
	for ruleId, rule := range ruleList {
		if strings.Contains(rule, label) {
			fwdRules[ruleId] = rule
		}
	}

	return fwdRules, nil
}

func DeleteForwardById(ruleId int) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	// delete rule
	err = ipt.Delete(FwdTable, FwdChain, strconv.Itoa(ruleId))
	if err != nil {
		return fmt.Errorf("failed deleting rule n. %d\nerr: %v", ruleId, err)
	}
	return nil
}

func DeleteForwardByRule(iface string, proto string, dport int, saddr string, sport int) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	// TODO: create function to return []string with packed rule, passing iface, proto, etc as arguments.
	ruleSpec := []string{
		"-i", iface,
		"-p", proto,
		"-m", proto,
		"--dport", strconv.Itoa(dport),
		"-m", "comment", "--comment", "fwdctl",
		"-j", FwdTarget,
		"--to-destination", saddr + ":" + strconv.Itoa(sport),
	}

	err = ipt.Delete(FwdTable, FwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("failed deleting rule: '%s'\n err: %v", ruleSpec, err)
	}
	return nil
}

func DeleteAllForwards() error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	ruleList, err := ipt.List(FwdTable, FwdChain)
	if err != nil {
		return fmt.Errorf("failed listing rules: %v", err)
	}

	// check listed rules are tagged with custom tag
	fwdRules := make(map[int]string)
	for ruleId, rule := range ruleList {
		if strings.Contains(rule, label) {
			fwdRules[ruleId] = rule
		}
	}

	for _, rule := range fwdRules {
		r, err := extractor.ExtractRuleInfo(rule)
		if err != nil {
			return fmt.Errorf("error extracting rule info: %v", err)
		}
		ruleSpec := []string{
			"-i", r[1],
			"-p", r[2],
			"-m", r[2],
			"--dport", r[3],
			"-m", "comment", "--comment", "fwdctl",
			"-j", FwdTarget,
			"--to-destination", r[4] + ":" + r[5],
		}
		err = ipt.Delete(FwdTable, FwdChain, ruleSpec...)
		if err != nil {
			return fmt.Errorf("error deleting rule: %v", err)
		}
	}
	return nil
}
