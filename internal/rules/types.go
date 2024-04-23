package rules

import (
	"github.com/alegrey91/fwdctl/pkg/iptables"
)

// RulesFile keep track of listed rules
// the result looks like so:
// rules:
//   - iface: eth0
//     proto: tcp
//     dport: 3000
//     saddr: 192.168.122.43
//     sport: 22
//   - iface: eth0
//     ...
type RuleSet struct {
	//Rules []iptables.Rule `yaml:"rules"`
	Rules map[string]iptables.Rule `yaml:"rules"`
}

// Struct to support creation of RuleSet
type supportRuleSet struct {
	Rules []iptables.Rule `yaml:"rules"`
}
