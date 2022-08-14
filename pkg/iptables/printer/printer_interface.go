package printer

import (
	"fmt"
	"strconv"
	"strings"
)

type Printer interface {
	PrintResult(ruleList []string) error
}

func NewPrinter(printFormat string) Printer {
	switch printFormat {
	case "table":
		return NewTable()
	default:
		return NewTable()
	}
}

func extractRuleInfo(rule string) ([]string, error) {
	// extract rules info:
	// -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	// result:
	// iface: eth0, protocol: tcp, dport: 3000, saddr: 192.168.199.105, sport: 80
	ruleSplit := strings.Split(rule, " ")
	fmt.Println(ruleSplit)
	ruleInfo := make([]string, 6)
	for id, arg := range ruleSplit {
		fmt.Println(arg)
		switch arg {
		case "-i":
			ruleInfo[1] = ruleSplit[id+1]
		case "-p":
			ruleInfo[2] = ruleSplit[id+1]
		case "--dport":
			ruleInfo[3] = ruleSplit[id+1]
		case "--to-destination":
			ruleInfo[4] = strings.Split(ruleSplit[id+1], ":")[0]
			ruleInfo[5] = strings.Split(ruleSplit[id+1], ":")[1]
		default:
			ruleInfo[0] = strconv.Itoa(id + 1)
		}
	}
	//fmt.Println("rule info: %v", ruleInfo)
	//for _, elem := range ruleInfo {
	//	fmt.Println("elem: %s", elem)
	//	if elem == "" {
	//		return []string{}, fmt.Errorf("unable to retrieve elements from rule")
	//	}
	//}
	//fmt.Println("rule info: %v", ruleInfo)
	return ruleInfo, nil
}
