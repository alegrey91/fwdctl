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
	"github.com/alegrey91/fwdctl/internal/template"
	st "github.com/alegrey91/fwdctl/internal/template/systemd_template"
	"github.com/spf13/cobra"
)

var installationPath string
var serviceType string

// generateSystemdCmd represents the generateSystemd command
var generateSystemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "Generates systemd service file",
	Long: `Generates systemd service file to run fwdctl at boot
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		systemd, err := st.NewSystemdService(serviceType, installationPath, c.RulesFile)
		if err != nil {
			return fmt.Errorf("cannot create systemd service: %v", err)
		}

		if err = template.GenerateTemplate(systemd, outputFile); err != nil {
			return fmt.Errorf("generating templated file: %v", err)
		}
		return nil
	},
}

func init() {
	generateCmd.AddCommand(generateSystemdCmd)

	generateSystemdCmd.Flags().StringVarP(&installationPath, "installation-path", "p", "/usr/local/bin", "fwdctl installation path")
	generateSystemdCmd.Flags().StringVarP(&c.RulesFile, "file", "f", "rules.yml", "rules file path")
	generateSystemdCmd.Flags().StringVarP(&serviceType, "type", "t", "oneshot", "systemd service type [oneshot, fork]")

	generateSystemdCmd.PersistentFlags().StringVarP(&outputFile, "output-path", "O", "", "output path")
	_ = generateSystemdCmd.MarkPersistentFlagRequired("output-path")
}
