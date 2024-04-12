package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var validateBlockCmd = &cobra.Command{
	Use:   "validateBlock channelID blockID",
	Short: "validate block by blockID and channel",
	Args:  cobra.ExactArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		channelID := args[0]
		if len(channelID) == 0 {
			FatalError("channelID is empty", nil)
		}

		blockID := args[1]
		if len(blockID) == 0 {
			FatalError("blockID is empty", nil)
		}
		validateBlock(channelID, blockID)
	},
}

func validateBlock(channelID string, blockID string) {
	logger.Info("validateBlock",
		zap.Any("channel", channelID),
		zap.Any("blockID", blockID),
	)
	file := channelID + "_" + blockID + ".block"
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		FatalError("Failed to WriteFile", err)
	}
	blk, err := protoutil.UnmarshalBlock(bytes)
	if err != nil {
		FatalError("Failed to Marshal", err)
	}
	if blk == nil || blk.Data == nil || len(blk.Data.Data) == 0 {
		FatalError("bad block", errors.New("invalid block"))
	}

	txs := blk.Data.Data
	for index, data := range txs {
		logger.Info("tx", zap.Any("tx", index))
		parseTx(data)
	}
}

func parseTx(envBytes []byte) {
	processedTransaction := peer.ProcessedTransaction{}
	err := proto.Unmarshal(envBytes, &processedTransaction)
	if err != nil {
		FatalError("unmarshal of transaction proposal processedTransaction failed", err)
	}
	logger.Info("processedTransaction.ValidationCode", zap.Any("ValidationCode", processedTransaction.ValidationCode))

	txID, err := protoutil.GetOrComputeTxIDFromEnvelope(envBytes)
	if err != nil {
		FatalError("error GetOrComputeTxIDFromEnvelope", err)
	}
	logger.Info("txID", zap.Any("txID", txID))

	transaction, err := protoutil.UnmarshalTransaction(envBytes)
	if err != nil {
		FatalError("error UnmarshalTransaction", err)
	}

	for i, transactionAction := range transaction.GetActions() {
		logger.Info("transactionAction ", zap.Any("i", i))
		rwSetByTransactionAction(transactionAction)
	}
}

func rwSetByTransactionAction(transactionAction *peer.TransactionAction) {
	_, ca, err := GetPayloads(transactionAction)
	if err != nil {
		return
	}

	chaincodeEvent := &peer.ChaincodeEvent{}
	if err = proto.Unmarshal(ca.Events, chaincodeEvent); err != nil {
		return
	}

	logger.Info("chaincodeEvent.EventName", zap.Any("EventName", chaincodeEvent.EventName))
}

func rwSetByChaincodeAction(chaincodeAction *peer.ChaincodeAction) {
	if chaincodeAction.Response.Status != http.StatusOK {
		return
	}

	chaincodeEvents, err := protoutil.UnmarshalChaincodeEvents(chaincodeAction.GetEvents())
	if err != nil {
		FatalError("chaincode events", err)
	}

	logger.Info("chaincodeEvents.EventName", zap.Any("chaincodeEvents.EventName", chaincodeEvents.EventName))

	txRWSet := &rwsetutil.TxRwSet{}
	err = txRWSet.FromProtoBytes(chaincodeAction.Results)
	if err != nil {
		FatalError("txRWSet From Proto Bytes", err)
	}
	ParseTxRWSet(txRWSet.NsRwSets)
}

func ParseTxRWSet(nsRwSet []*rwsetutil.NsRwSet) {
	for _, set := range nsRwSet {
		printReads(set)
		printWrites(set)
		printRangeQueriesInfo(set)
		printMetadataWrites(set)
		printCollectionName(set)
	}
}

func printReads(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.Reads) == 0 {
		return
	}
	logger.Info("KvRwSet.Reads", zap.Any("Reads", set.KvRwSet.Writes))
}

func printWrites(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.Writes) == 0 {
		return
	}
	logger.Info("KvRwSet.Writes", zap.Any("Writes", set.KvRwSet.Writes))
}

func printRangeQueriesInfo(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.RangeQueriesInfo) == 0 {
		return
	}
	logger.Info("KvRwSet.RangeQueriesInfo", zap.Any("RangeQueriesInfo", set.KvRwSet.RangeQueriesInfo))
}

func printMetadataWrites(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.MetadataWrites) == 0 {
		return
	}
	logger.Info("KvRwSet.MetadataWrites", zap.Any("MetadataWrites", set.KvRwSet.MetadataWrites))
}

func printCollectionName(set *rwsetutil.NsRwSet) {
	if len(set.CollHashedRwSets) == 0 {
		return
	}
	logger.Info("set.CollectionName", zap.Any("CollHashedRwSets", set.CollHashedRwSets))
}

// GetPayloads gets the underlying payload objects in a TransactionAction.
func GetPayloads(txActions *peer.TransactionAction) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	ccPayload, err := GetChaincodeActionPayload(txActions.Payload)
	if err != nil {
		return nil, nil, err
	}

	pRespPayload, err := GetProposalResponsePayload(
		ccPayload.ChaincodeProposalPayload,
	)
	if err != nil {
		return nil, nil, err
	}

	logger.Info("pRespPayload", zap.Any("pRespPayload", pRespPayload))

	if pRespPayload.Extension == nil {
		return nil, nil, errors.New("extension is nil for txActions.Payload.ChaincodeProposalPayload")
	}

	respPayload, err := GetChaincodeAction(pRespPayload.Extension)
	if err != nil {
		return nil, nil, err
	}

	payload, err := GetChaincodeActionPayload(respPayload.Results)
	if err != nil {
		return nil, nil, err
	}
	responsePayload, err := GetProposalResponsePayload(payload.GetChaincodeProposalPayload())
	if err != nil {
		return nil, nil, err
	}

	logger.Info("responsePayload",
		zap.Any("responsePayload.String()", responsePayload.String()),
	)
	logger.Info("payload",
		zap.Any("payload.ChaincodeProposalPayload", payload.ChaincodeProposalPayload),
	)

	return ccPayload, respPayload, nil
}

// GetChaincodeActionPayload Get ChaincodeActionPayload from bytes.
func GetChaincodeActionPayload(
	capBytes []byte,
) (*peer.ChaincodeActionPayload, error) {
	c := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(capBytes, c)
	if err != nil {
		return c, fmt.Errorf("error unmarshaling ChaincodeActionPayload %w", err)
	}
	return c, nil
}

// GetProposalResponsePayload gets the proposal response payload.
func GetProposalResponsePayload(
	prpBytes []byte,
) (*peer.ProposalResponsePayload, error) {
	prp := &peer.ProposalResponsePayload{}
	err := proto.Unmarshal(prpBytes, prp)
	if err != nil {
		return prp, fmt.Errorf("error unmarshaling ProposalResponsePayload %w", err)
	}
	return prp, nil
}

// GetChaincodeAction gets the ChaincodeAction given chaicnode action bytes.
func GetChaincodeAction(caBytes []byte) (*peer.ChaincodeAction, error) {
	chaincodeAction := &peer.ChaincodeAction{}
	err := proto.Unmarshal(caBytes, chaincodeAction)
	if err != nil {
		return chaincodeAction, fmt.Errorf("error unmarshaling ChaincodeAction %w", err)
	}
	return chaincodeAction, nil
}
