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
	"os"

	//"os"

	c "github.com/alegrey91/fwdctl/internal/constants"
	"github.com/alegrey91/fwdctl/internal/rules"
	ipt "github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/spf13/cobra"
)

var (
	ruleId int
	file   string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm"},
	Short:   "Delete forward",
	Long: `Delete forward by passing a rule file or rule id.
`,
	Example: c.ProgramName + " delete -n 2",
	Run: func(cmd *cobra.Command, args []string) {

		// Delete rule number
		if cmd.Flags().Lookup("id").Changed {
			err := ipt.DeleteForwardById(ruleId)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		// Loop over file content and delete rule one-by-one.
		if cmd.Flags().Lookup("file").Changed {
			deleteFromFile(file)
			return
		}

		if cmd.Flags().Lookup("all").Changed {
			err := ipt.DeleteAllForwards()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		deleteFromFile(file)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().IntVarP(&ruleId, "id", "n", 0, "delete rules through ID")
	deleteCmd.Flags().StringVarP(&file, "file", "f", "rules.yml", "delete rules through file")
	deleteCmd.Flags().BoolP("all", "a", false, "delete all rules")
	deleteCmd.MarkFlagsMutuallyExclusive("id", "file", "all")
}

func deleteFromFile(file string) {
	rulesFile, err := rules.NewRuleSetFromFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, rule := range rulesFile.Rules {
		err := ipt.DeleteForwardByRule(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
