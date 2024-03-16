package iptables

var (
	FwdTable  = "nat"
	FwdChain  = "PREROUTING"
	FwdTarget = "DNAT"
)
