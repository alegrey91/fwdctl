package iptables

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alegrey91/fwdctl/internal/extractor"
	"github.com/coreos/go-iptables/iptables"
)

var (
	label string = "fwdctl"
)

type IPTablesInstance struct {
	*iptables.IPTables
}

func NewIPTablesInstance() (*IPTablesInstance, error) {
	ipt := IPTablesInstance{}
	iptables, err := getIPTablesInstance()
	if err != nil {
		return nil, fmt.Errorf("failed: %v", err)
	}
	ipt.IPTables = iptables
	return &ipt, nil
}

func (ipt *IPTablesInstance) ValidateForward(iface string, proto string, dport int, saddr string, sport int) error {
	return validate(iface, proto, dport, saddr, sport)
}

func (ipt *IPTablesInstance) CreateForward(iface string, proto string, dport int, saddr string, sport int) error {
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

func (ipt *IPTablesInstance) ListForward(outputFormat string) (map[int]string, error) {
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

func (ipt *IPTablesInstance) DeleteForwardById(ruleId int) error {
	// delete rule
	err := ipt.Delete(FwdTable, FwdChain, strconv.Itoa(ruleId))
	if err != nil {
		return fmt.Errorf("failed deleting rule n. %d\nerr: %v", ruleId, err)
	}
	return nil
}

func (ipt *IPTablesInstance) DeleteForwardByRule(iface string, proto string, dport int, saddr string, sport int) error {
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

	err := ipt.Delete(FwdTable, FwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("failed deleting rule: '%s'\n err: %v", ruleSpec, err)
	}
	return nil
}

func (ipt *IPTablesInstance) DeleteAllForwards() error {
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
