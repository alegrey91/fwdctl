package iptables

import (
	"fmt"
	"strconv"
	"strings"
)

type Rule struct {
	Iface   string `default:"lo"`
	Proto   string `default:"tcp"`
	Dport   int
	Saddr   string
	Sport   int
	Comment string
}

func NewRule(iface string, proto string, dport int, saddr string, sport int) *Rule {
	return &Rule{
		Iface:   iface,
		Proto:   proto,
		Dport:   dport,
		Saddr:   saddr,
		Sport:   sport,
		Comment: label,
	}
}

// ExtractRuleInfo extract forward information from rule
// if it matches the requirements.
// Returns the Rule struct and error
func ExtractRuleInfo(rawRule string) (*Rule, error) {
	// extract rules info:
	// -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	// result:
	// Rule{Iface: eth0, Proto: tcp, Dport: 3000, Saddr: 192.168.199.105, Sport: 80}
	ruleSplit := strings.Split(rawRule, " ")
	rule := &Rule{}
	for id, arg := range ruleSplit {
		switch arg {
		case "-i":
			rule.Iface = ruleSplit[id+1]
		case "-p":
			rule.Proto = ruleSplit[id+1]
		case "--dport":
			dport, err := strconv.Atoi(ruleSplit[id+1])
			if err != nil {
				return nil, fmt.Errorf("error converting string '%s' to int: %v", ruleSplit[id+1], err)
			}
			rule.Dport = dport
		case "--to-destination":
			rule.Saddr = strings.Split(ruleSplit[id+1], ":")[0]
			sport, err := strconv.Atoi(strings.Split(ruleSplit[id+1], ":")[1])
			if err != nil {
				return nil, fmt.Errorf("error converting string '%s' to int: %v", ruleSplit[id+1], err)
			}
			rule.Sport = sport
		}
	}

	return rule, nil
}

// String returns a list of string that compose the iptables rule.
// Eg: -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -m comment --comment fwdctl -j DNAT --to-destination 192.168.199.105:80
func (rule *Rule) String() []string {
	return []string{
		"-i", rule.Iface,
		"-p", rule.Proto,
		"-m", rule.Proto,
		"--dport", fmt.Sprintf("%d", rule.Dport),
		"-m", "comment", "--comment", label,
		"-j", FwdTarget,
		"--to-destination", rule.Saddr + ":" + fmt.Sprintf("%d", rule.Sport),
	}
}
