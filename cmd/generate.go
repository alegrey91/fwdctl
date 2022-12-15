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
)

var outputFile string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Aliases:    []string{"gen"},
	Short: "generates templated files",
	Long: `generates templated file for fwdtcl
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
		    err := cmd.Help()
			if err != nil {
				fmt.Println(err)
			}
		    os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "output path")
	err := generateCmd.MarkPersistentFlagRequired("output")
	if err != nil {
		fmt.Println(err.Error())
        return	
	}
}
