package printer

import (
	"fmt"
	"strconv"

	yaml "gopkg.in/yaml.v3"

	"github.com/alegrey91/fwdctl/internal/rules"
	"github.com/alegrey91/fwdctl/pkg/iptables"
)

type Yaml struct {
}

func NewYaml() *Yaml {
	return &Yaml{}
}

func (y *Yaml) PrintResult(ruleList map[int]string) error {
	rules := rules.NewRuleSet()
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
		rules.Add(rule)
	}
	val, err := yaml.Marshal(rules.Rules)
	if err != nil {
		return err
	}

	fmt.Println(string(val))
	return nil
}
