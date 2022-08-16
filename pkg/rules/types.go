package rules

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
type RulesFile struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Iface string `yaml:"iface"`
	Proto string `yaml:"proto"`
	Dport int    `yaml:"dport"`
	Saddr string `yaml:"saddr"`
	Sport int    `yaml:"sport"`
}
