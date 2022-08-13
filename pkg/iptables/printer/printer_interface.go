package printer

type Printer interface {
	PrintResult(ruleList [][]string)
}

func ExtractRuleInfo(rules []string) [][]string {
	// extract rules info:
	// -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	// result:
	// iface: eth0, protocol: tcp, dport: 3000, saddr: 192.168.199.105, sport: 80
}