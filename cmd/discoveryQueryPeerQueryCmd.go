package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var discoveryQueryPeerQueryCmd = &cobra.Command{ //nolint:unused
	Use:   "discoveryQueryPeerQuery",
	Short: "discovery query peer Query",
	Args:  cobra.ExactArgs(3), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()
		fmt.Println("command not implement yet")
	},
}
