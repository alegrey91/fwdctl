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
	"github.com/spf13/cobra"
)

var (
	ruleId int
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete forward by passing Id number",
	Long: `Delete forward by passing Id number.
The Id number is retrieved using the command:
sudo iptables -t nat -L PREROUTING -n --line-number
`,
	Example: c.ProgramName + " delete -n 2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ruleId)
		err := ipt.DeleteForward(ruleId)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().IntVarP(&ruleId, "rule-id", "n", 0, "rule number")
	deleteCmd.MarkFlagRequired("rule-id")
}
