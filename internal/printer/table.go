package printer

import (
	"fmt"
	"os"

	"github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/olekukonko/tablewriter"
)

type Table struct {
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) PrintResult(ruleList map[int]string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"number", "interface", "protocol", "external port", "internal ip", "internal port"})
	for ruleId, rule := range ruleList {
		tabRule, err := iptables.ExtractRuleInfo(rule)
		if err != nil {
			continue
		}
		tabRow := []string{
			fmt.Sprintf("%d", ruleId),
			tabRule.Iface,
			tabRule.Proto,
			fmt.Sprintf("%d", tabRule.Dport),
			tabRule.Saddr,
			fmt.Sprintf("%d", tabRule.Sport),
		}
		table.Append(tabRow)
	}
	table.Render()
	return nil
}
