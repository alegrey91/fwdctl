package iptables

type Rule struct {
	Iface string
	Proto string
	Dport int
	Saddr string
	Sport int
}