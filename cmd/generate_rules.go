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

	"github.com/spf13/cobra"

	"github.com/alegrey91/fwdctl/internal/template"
	rt "github.com/alegrey91/fwdctl/internal/template/rules_template"
)

// generateRulesCmd represents the generateRules command
var generateRulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Generates empty rules file",
	Long: `Generates empty rules file
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rules := rt.NewRules()
		if err := template.GenerateTemplate(rules, outputFile); err != nil {
			return fmt.Errorf("generating template: %w", err)
		}
		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateRulesCmd)

	generateRulesCmd.PersistentFlags().StringVarP(&outputFile, "output-path", "O", "", "output path")
	_ = generateRulesCmd.MarkPersistentFlagRequired("output-path")
}
