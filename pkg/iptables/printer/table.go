package printer

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func PrintResult(ruleList []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"number", "external port", "internal IP", "internal port", "protocol"})
	for _, rule := range ruleList {
		//table.Append(rule)
		
	}
	table.Render()
}
