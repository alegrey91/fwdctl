package printer

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alegrey91/fwdctl/internal/rules"
	"github.com/alegrey91/fwdctl/pkg/iptables"
)

type Json struct {
}

func NewJson() *Json {
	return &Json{}
}

func (j *Json) PrintResult(ruleList []string) error {
	rules := rules.RulesFile{}
	for _, rule := range ruleList {
		jsonRule, err := extractRuleInfo(rule)
		if err != nil {
			continue
		}

		dport, err := strconv.Atoi(jsonRule[3])
		if err != nil {
			return err
		}
		sport, err := strconv.Atoi(jsonRule[5])
		if err != nil {
			return err
		}
		rule := iptables.Rule{
			Iface: jsonRule[1],
			Proto: jsonRule[2],
			Dport: dport,
			Saddr: jsonRule[4],
			Sport: sport,
		}
		rules.Rules = append(rules.Rules, rule)
	}
	val, err := json.MarshalIndent(rules.Rules, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(val))
	return nil
}
