package cmd

import (
	"strconv"

	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/spf13/cobra"
)

type TxStatInfo struct {
	TxType         cb.HeaderType
	ChaincodeID    *peer.ChaincodeID
	NumCollections []*rwsetutil.NsRwSet
}

var txCmd = &cobra.Command{
	Use:   "tx channelID transactionID peerUrl",
	Short: "get transaction and block by transactionId chaincode",
	Args:  cobra.ExactArgs(3), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID := args[0]
		if len(channelID) == 0 {
			FatalError("channelID is empty", nil)
		}

		transactionID := args[1]
		if len(transactionID) == 0 {
			FatalError("transactionID is empty", nil)
		}

		peer := args[2]
		if len(peer) == 0 {
			FatalError("peer is empty", nil)
		}

		block, err := HlfClient.QueryBlockByTxID(channelID, transactionID, peer)
		if err != nil {
			FatalError("QueryBlockByTxID", err)
		}

		blockID := strconv.FormatUint(block.GetHeader().GetNumber(), 10)
		err = saveBlock(channelID, blockID, block)
		if err != nil {
			FatalError("failed to save block", err)
		}
	},
}
