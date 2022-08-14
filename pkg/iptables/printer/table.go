package printer

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) PrintResult(ruleList []string) error {
	table := tablewriter.NewWriter(os.Stdout) 
	table.SetHeader([]string{"number", "interface", "protocol", "external port", "internal ip", "internal port"})
	for _, rule := range ruleList {
		fmt.Println(rule)
		tabRule, err := extractRuleInfo(rule)
		if err != nil {
			return err
		}
		table.Append(tabRule)
	}
	table.Render()
	return nil
}
