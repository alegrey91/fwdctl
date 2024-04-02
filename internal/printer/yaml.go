package printer

import (
	"fmt"

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
		jsonRule, err := iptables.ExtractRuleInfo(rule)
		if err != nil {
			continue
		}

		rules.Add(*jsonRule)
	}
	val, err := yaml.Marshal(rules.Rules)
	if err != nil {
		return err
	}

	fmt.Println(string(val))
	return nil
}
