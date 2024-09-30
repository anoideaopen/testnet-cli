package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var channelHeightCmd = &cobra.Command{
	Use:   "channelHeight channelID peer",
	Short: "get channel height on peer",
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

		blockchainInfoResponse, err := HlfClient.GetPeerInfo(channelID, peer)
		if err != nil {
			FatalError("Failed to get peer info", err)
		}
		fmt.Println(blockchainInfoResponse.BCI.GetHeight())
	},
}
