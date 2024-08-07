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
	iptables "github.com/alegrey91/fwdctl/pkg/iptables"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ipt, err := iptables.NewIPTablesInstance()
		if err != nil {
			return fmt.Errorf("unable to get iptables instance: %v", err)
		}
		// Delete rule number
		if cmd.Flags().Lookup("id").Changed {
			if err := ipt.DeleteForwardById(ruleId); err != nil {
				return fmt.Errorf("delete forward by ID: %v", err)
			}
			return nil
		}

		// Loop over file content and delete rule one-by-one.
		if cmd.Flags().Lookup("file").Changed {
			if err := deleteFromFile(ipt, file); err != nil {
				return fmt.Errorf("delete from file: %v", err)
			}
			return nil
		}

		if cmd.Flags().Lookup("all").Changed {
			if err := ipt.DeleteAllForwards();err != nil {
				return fmt.Errorf("delete all forwards: %v", err)
			}
			return nil
		}

		if err = deleteFromFile(ipt, file);err != nil {
			return fmt.Errorf("delete from file: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().IntVarP(&ruleId, "id", "n", 0, "delete rules through ID")
	deleteCmd.Flags().StringVarP(&file, "file", "f", "rules.yml", "delete rules through file")
	deleteCmd.Flags().BoolP("all", "a", false, "delete all rules")
	deleteCmd.MarkFlagsMutuallyExclusive("id", "file", "all")
}

func deleteFromFile(ipt *iptables.IPTablesInstance, file string) error {
	rulesContent, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	rulesFile, err := rules.NewRuleSetFromFile(rulesContent)
	if err != nil {
		return fmt.Errorf("error instantiating ruleset from file: %v", err)
	}
	for _, rule := range rulesFile.Rules {
		err := ipt.DeleteForwardByRule(&rule)
		if err != nil {
			return fmt.Errorf("error deleting rule [%s %s %d %s %d]: %v", rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport, err)
		}
	}
	return nil
}
