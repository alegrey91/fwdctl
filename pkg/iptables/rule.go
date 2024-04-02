package iptables

import "fmt"

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
