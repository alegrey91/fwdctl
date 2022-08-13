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
	iface string
	proto string
	dport int
	saddr string
	sport int
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:        "create",
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "Create forward using IPTables util",
	Long: `Create forward rule using IPTables util under the hood.
This is really useful in case you need to forward
the traffic from an internal virtual machine inside
your hypervisor, to external.

   +----------------------------+
   |              +-----------+ |
   |              |           | |
   |        +-----+:80  VM    | |
   |        |     |           | |
   =:3000<--+     +-----------+ |
   |         Hypervisor         |
   +----------------------------+
`,
	Example: c.ProgramName + " create -d 3000 -s 192.168.199.105 -p 80",
	Run: func(cmd *cobra.Command, args []string) {
		err := ipt.CreateForward(iface, proto, dport, saddr, sport)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&iface, "interface", "i", "eth0", "interface name")
	createCmd.Flags().StringVarP(&proto, "proto", "P", "tcp", "protocol")

	createCmd.Flags().IntVarP(&dport, "destination-port", "d", 0, "destination port")
	createCmd.MarkFlagRequired("destination-port")

	createCmd.Flags().StringVarP(&saddr, "source-address", "s", "", "source address")
	createCmd.MarkFlagRequired("source-address")

	createCmd.Flags().IntVarP(&sport, "source-port", "p", 0, "source port")
	createCmd.MarkFlagRequired("source-port")
}
