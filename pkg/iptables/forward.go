package iptables

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alegrey91/fwdctl/pkg/iptables/printer"
)

func CreateForward(iface string, proto string, dport int, saddr string, sport int) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	rule := Rule{
		Iface: iface,
		Proto: proto,
		Dport: dport,
		Saddr: saddr,
		Sport: sport,
	}

	// example rule:
	// iptables -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	ruleSpec := []string{
		"-i", rule.Iface,
		"-p", rule.Proto,
		"-m", rule.Proto,
		"--dport", strconv.Itoa(rule.Dport),
		"-j", fwdTarget,
		"--to-destination", rule.Saddr + ":" + strconv.Itoa(rule.Sport),
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
	ruleExists, err := ipt.Exists(fwdTable, fwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if ruleExists {
		return fmt.Errorf("rule already exists")
	}

	// apply provided rule
	err = ipt.AppendUnique(fwdTable, fwdChain, ruleSpec...)
	if err != nil {
		return fmt.Errorf("rule failed: %v", err)
	}
	return nil
}

func ListForward(outputFormat string) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	//ruleList, err := ipt.ListWithCounters(fwdTable, fwdChain)
	ruleList, err := ipt.List(fwdTable, fwdChain)
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}
	
	p := printer.NewPrinter(outputFormat)
	err = p.PrintResult(ruleList)
	if err != nil {
		return fmt.Errorf("failed printing results: %v", err)
	}
	return nil
}

func DeleteForward(ruleId int) error {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return fmt.Errorf("failed: %v", err)
	}

	// retrieve rule using Id number
	// (sudo iptables -t nat -L PREROUTING -n --line-numbers)
	rule, err := ipt.ListById(fwdTable, fwdChain, ruleId)
	if err != nil {
		return fmt.Errorf("unable to retrieve rule with ID: %d", ruleId)
	}

	// cleaning rule (removing "-A PREROUTING", "-c 0 0", ...)
	ruleSplit := strings.Split(rule, " ")
	ruleSplit = append(ruleSplit[2:10], ruleSplit[13:]...)

	// delete rule
	err = ipt.Delete(fwdTable, fwdChain, ruleSplit...)
	if err != nil {
		return fmt.Errorf("failed deleting rule n. %d\n err: %v", ruleId, err)
	}
	return nil
}
