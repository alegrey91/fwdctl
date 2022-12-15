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

	"github.com/alegrey91/fwdctl/pkg/template"
	rt "github.com/alegrey91/fwdctl/pkg/template/rules_template"
	"github.com/spf13/cobra"
)

// generateRulesCmd represents the generateRules command
var generateRulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "generates empty rules file",
	Long: `generates empty rules file
`,
	Run: func(cmd *cobra.Command, args []string) {
		rules := rt.NewRules()
		err := template.GenerateTemplate(rules, outputFile)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateRulesCmd)
}
