package printer

import (
	"encoding/json"
	"fmt"

	"github.com/alegrey91/fwdctl/internal/rules"
	"github.com/alegrey91/fwdctl/pkg/iptables"
)

type Json struct {
}

func NewJson() *Json {
	return &Json{}
}

func (j *Json) PrintResult(ruleList map[int]string) error {
	rules := rules.NewRuleSet()
	for _, rule := range ruleList {
		jsonRule, err := iptables.ExtractRuleInfo(rule)
		if err != nil {
			continue
		}

		rules.Add(*jsonRule)
	}
	val, err := json.MarshalIndent(rules.Rules, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(val))
	return nil
}
