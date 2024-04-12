package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var chaincodeVersionCmd = &cobra.Command{
	Use:   "v channelID",
	Short: "v - version chaincode",
	Args:  cobra.ExactArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID := args[0]
		if len(channelID) == 0 {
			FatalError("channelID is empty", nil)
		}

		peer := args[1]
		if len(peer) == 0 {
			FatalError("peer is empty", nil)
		}

		version, err := HlfClient.ChaincodeVersion(channelID, peer)
		if err != nil {
			FatalError("Failed to create new channel client", err)
		}
		fmt.Println(version)
	},
}
