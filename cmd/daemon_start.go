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

	c "github.com/alegrey91/fwdctl/internal/constants"
	"github.com/alegrey91/fwdctl/internal/daemon"
	iptables "github.com/alegrey91/fwdctl/pkg/iptables"
	"github.com/spf13/cobra"
)

// daemonStartCmd represents the daemon command
var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start fwdctl daemon",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ipt, err := iptables.NewIPTablesInstance()
		if err != nil {
			fmt.Printf("unable to get iptables instance: %v\n", err)
			os.Exit(1)
		}
		rulesFile, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("unable to read from flag: %v", err)
			os.Exit(1)
		}
		if res := daemon.Start(ipt, rulesFile); res != 0 {
			os.Exit(1)
		}
	},
}

func init() {
	daemonCmd.AddCommand(daemonStartCmd)
	daemonStartCmd.Flags().StringVarP(&c.RulesFile, "file", "f", "rules.yml", "rules file")
}
