package iptables

type Rule struct {
	Iface string `default:"lo"`
	Proto string `default:"tcp"`
	Dport int
	Saddr string
	Sport int
}