package printer

import (
	"fmt"
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

// extractRuleInfo extract forward information from rule
// if it matches the requirements
func extractRuleInfo(rule string) ([]string, error) {
	// extract rules info:
	// -t nat -A PREROUTING -i eth0 -p tcp -m tcp --dport 3000 -j DNAT --to-destination 192.168.199.105:80
	// result:
	// slice ruleInfo: [ , eth0, tcp, 3000, 192.168.199.105, 80]
	ruleSplit := strings.Split(rule, " ")
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
		}
	}

	for i := 1; i < len(ruleInfo); i++ {
		if ruleInfo[i] == "" {
			return []string{}, fmt.Errorf("unable to retrieve elements from rule")
		}
	}
	return ruleInfo, nil
}
