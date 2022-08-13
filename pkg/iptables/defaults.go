package iptables

var (
	fwdTable  = "nat"
	fwdChain  = "PREROUTING"
	fwdTarget = "DNAT"
)