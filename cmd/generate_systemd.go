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
	"github.com/spf13/cobra"
)

var installationPathGen string
var rulesFileGen string

// generateSystemdCmd represents the generateSystemd command
var generateSystemdCmd = &cobra.Command{
	Use:   "systemd",
	Short: "generates systemd service file",
	Long: `generates systemd service file to run fwdctl at boot
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := template.GetSystemdTemplate(outputFile, installationPathGen, rulesFileGen)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateSystemdCmd)
	

	generateSystemdCmd.Flags().StringVarP(&installationPathGen, "installation-path", "p", "/usr/local/bin", "fwdctl installation path")
	generateSystemdCmd.Flags().StringVarP(&rulesFileGen, "rules-file", "r", "~/rules.yml", "rules file")
}
