package iptables

import (
	"fmt"
	"strconv"
	"strings"

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

func (ipt *IPTablesInstance) ValidateForward(rule *Rule) error {
	return validate(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
}

func (ipt *IPTablesInstance) CreateForward(rule *Rule) error {
	// example rule:
	// iptables -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80

	// check if input interface exists on the system
	ifaceExits, err := interfaceExists(rule.Iface)
	if err != nil {
		return fmt.Errorf("error reading interfaces: %v", err)
	}
	if !ifaceExits {
		return fmt.Errorf("interface %s does not exists", rule.Iface)
	}

	// check if provided rule already exists
	ruleExists, err := ipt.Exists(FwdTable, FwdChain, rule.String()...)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if ruleExists {
		return fmt.Errorf("rule already exists")
	}

	// apply provided rule
	err = ipt.AppendUnique(FwdTable, FwdChain, rule.String()...)
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

func (ipt *IPTablesInstance) DeleteForwardByRule(rule *Rule) error {
	// TODO: create function to return []string with packed rule, passing iface, proto, etc as arguments.
	err := ipt.Delete(FwdTable, FwdChain, rule.String()...)
	if err != nil {
		return fmt.Errorf("failed deleting rule: '%s'\n err: %v", rule.String(), err)
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
		r, err := ExtractRuleInfo(rule)
		if err != nil {
			return fmt.Errorf("error extracting rule info: %v", err)
		}

		err = ipt.Delete(FwdTable, FwdChain, r.String()...)
		if err != nil {
			return fmt.Errorf("error deleting rule: %v", err)
		}
	}
	return nil
}
