package printer

import (
	"os"
	"strconv"

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
		tabRule, err := extractRuleInfo(rule)
		if err != nil {
			continue
		}
		tabRule[0] = strconv.Itoa(ruleId)
		table.Append(tabRule)
	}
	table.Render()
	return nil
}
