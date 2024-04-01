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
	"sync"

	c "github.com/alegrey91/fwdctl/internal/constants"
	"github.com/alegrey91/fwdctl/internal/rules"
	iptables "github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:     "apply",
	Short:   "apply rules from file",
	Long:    `apply rules described in a configuration file`,
	Example: c.ProgramName + " apply --file rule.yml",
	Run: func(cmd *cobra.Command, args []string) {
		ruleSet, err := rules.NewRuleSetFromFile(c.RulesFile)
		if err != nil {
			fmt.Printf("unable to open rules file: %v\n", err)
			os.Exit(1)
		}
		ipt, err := iptables.NewIPTablesInstance()
		if err != nil {
			fmt.Printf("unable to get iptables instance: %v\n", err)
			os.Exit(1)
		}

		var wg sync.WaitGroup
		chErr := make(chan error, len(ruleSet.Rules))
		chLimit := make(chan int, 10)
		rulesFileIsValid := true

		for _, rule := range ruleSet.Rules {
			wg.Add(1)
			// add slot to buffered channel
			chLimit <- 1
			go func(rule iptables.Rule, wg *sync.WaitGroup, chErr chan error, chLimit chan int) {
				err := ipt.ValidateForward(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
				wg.Done()
				chErr <- err
				// free slot from buffered channel
				<-chLimit
			}(rule, &wg, chErr, chLimit)
		}
		go func() {
			wg.Wait()
			close(chErr)
		}()

		for err := range chErr {
			if err != nil {
				fmt.Printf("error validating rule: %v\n", err)
				os.Exit(1)
			}
		}

		if rulesFileIsValid {
			for ruleId, rule := range ruleSet.Rules {
				err = ipt.CreateForward(rule.Iface, rule.Proto, rule.Dport, rule.Saddr, rule.Sport)
				if err != nil {
					fmt.Printf("error applying rule (%s): %v\n", ruleId, err)
					os.Exit(1)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVarP(&c.RulesFile, "file", "f", "rules.yml", "rules file")
}
