/*
Copyright © 2022 Alessio Greggi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	c "github.com/alegrey91/fwdctl/internal/constants"
	ipt "github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/alegrey91/fwdctl/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	format string
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list forwards",
	Long:    `list forwards made with iptables`,
	Example: c.ProgramName + "list -o table",
	Run: func(cmd *cobra.Command, args []string) {
		ruleList, err := ipt.ListForward(format)
		if err != nil {
			fmt.Println(err)
			return
		}

		p := printer.NewPrinter(format)
		err = p.PrintResult(ruleList)
		if err != nil {
			fmt.Printf("failed printing results: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&format, "output", "o", "table", "output format [table]")
}
