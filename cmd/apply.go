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

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	c "github.com/alegrey91/fwdctl/internal/constants"
	"github.com/alegrey91/fwdctl/internal/rules"
	iptables "github.com/alegrey91/fwdctl/pkg/iptables"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:     "apply",
	Short:   "Apply rules from file",
	Long:    `Apply rules described in a configuration file`,
	Example: c.ProgramName + " apply --file rule.yml",
	RunE: func(cmd *cobra.Command, args []string) error {
		rulesContent, err := os.Open(c.RulesFile)
		if err != nil {
			return fmt.Errorf("opening file: %v", err)
		}
		ruleSet, err := rules.NewRuleSetFromFile(rulesContent)
		if err != nil {
			return fmt.Errorf("unable to open rules file: %v", err)
		}
		ipt, err := iptables.NewIPTablesInstance()
		if err != nil {
			return fmt.Errorf("unable to get iptables instance: %v", err)
		}

		g := new(errgroup.Group)
		rulesFileIsValid := true

		g.SetLimit(10)
		for _, rule := range ruleSet.Rules {
			r := &rule
			g.Go(func() error {
				err := ipt.ValidateForward(r)
				if err != nil {
					rulesFileIsValid = false
				}
				return err
			})
		}

		if err := g.Wait(); err != nil {
			return fmt.Errorf("validating rule: %v", err)
		}

		if rulesFileIsValid {
			for ruleId, rule := range ruleSet.Rules {
				if err := ipt.CreateForward(&rule); err != nil {
					return fmt.Errorf("applying rule (%s): %v", ruleId, err)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVarP(&c.RulesFile, "file", "f", "rules.yml", "rules file")
}
