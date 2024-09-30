package cmd

import (
	"errors"
	"os"
	"strconv"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var blockByIDCmd = &cobra.Command{
	Use:   "block channelID blockID peer",
	Short: "save block by blockID to same dir and print block info",
	Args:  cobra.ExactArgs(3), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID := args[0]
		if len(channelID) == 0 {
			FatalError("channelID is empty", nil)
		}

		blockID := args[1]
		if len(blockID) == 0 {
			FatalError("blockID is empty", nil)
		}

		peer := args[2]
		if len(peer) == 0 {
			FatalError("peer is empty", nil)
		}
		block, err := getBlock(channelID, blockID, peer)
		if err != nil {
			FatalError("failed to get block", err)
		}
		if block == nil {
			FatalError("failed to get block, block nil", errors.New("failed to get block, block nil"))
		}
		blockID = strconv.FormatUint(block.GetHeader().GetNumber(), 10)
		err = saveBlock(channelID, blockID, block)
		if err != nil {
			FatalError("failed to save block", err)
		}
		printBlock(block)
	},
}

func getBlock(channelID string, blockID string, peer string) (*cb.Block, error) {
	block, err := HlfClient.QueryBlock(channelID, blockID, peer)
	if err != nil {
		FatalError("Failed to QueryBlock", err)
		return nil, err
	}

	return block, nil
}

func saveBlock(channelID string, blockID string, block *cb.Block) error {
	b, err := proto.Marshal(block)
	if err != nil {
		return err
	}

	file := channelID + "_" + blockID + ".block"

	err = os.WriteFile(file, b, 0o600) //nolint:gomnd
	if err != nil {
		return err
	}

	return nil
}

func printBlock(block *cb.Block) {
	for txIndex, envBytes := range block.GetData().GetData() {
		logger.Info("txIndex", zap.Any("txIndex", txIndex))
		env, err := protoutil.GetEnvelopeFromBlock(envBytes)
		if err != nil {
			FatalError("GetEnvelopeFromBlock", err)
		}

		payload, err := protoutil.UnmarshalPayload(env.GetPayload())
		if err != nil {
			FatalError("UnmarshalPayload", err)
		}

		shdr, err := protoutil.UnmarshalSignatureHeader(payload.GetHeader().GetSignatureHeader())
		if err != nil {
			FatalError("UnmarshalSignatureHeader", err)
		}
		logger.Info("shdr.Creator")
		logger.Info(string(shdr.GetCreator()))

		serializedIdentity, err := protoutil.UnmarshalSerializedIdentity(shdr.GetCreator())
		if err != nil {
			FatalError("UnmarshalSerializedIdentity", err)
		}

		logger.Info("serializedIdentity.Mspid")
		logger.Info(serializedIdentity.GetMspid())

		logger.Info("serializedIdentity.IdBytes")
		logger.Info(string(serializedIdentity.GetIdBytes()))

		logger.Info("shdr.Nonce")
		logger.Info(string(shdr.GetNonce()))

		chdr, err := protoutil.UnmarshalChannelHeader(payload.GetHeader().GetChannelHeader())
		if err != nil {
			FatalError("UnmarshalChannelHeader", err)
		}

		logger.Info("chdr.TxId")
		logger.Info(chdr.GetTxId())

		logger.Info("payload.Data")
		logger.Info(string(payload.GetData()))
	}
}
