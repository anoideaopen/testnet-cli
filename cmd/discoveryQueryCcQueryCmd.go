package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var discoveryQueryCcQueryCmd = &cobra.Command{
	Use:   "discoveryQueryCcQuery",
	Short: "discovery Query Cc Query",
	Args:  cobra.ExactArgs(3), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()
		fmt.Println("command not implement yet")
	},
}
