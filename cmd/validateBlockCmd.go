package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var validateBlockCmd = &cobra.Command{ //nolint:unused
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

func validateBlock(channelID string, blockID string) { //nolint:unused
	logger.Info("validateBlock",
		zap.Any("channel", channelID),
		zap.Any("blockID", blockID),
	)
	file := channelID + "_" + blockID + ".block"
	bytes, err := os.ReadFile(file)
	if err != nil {
		FatalError("Failed to WriteFile", err)
	}
	blk, err := UnmarshalBlock(bytes)
	if err != nil {
		FatalError("Failed to Marshal", err)
	}
	if len(blk.GetData().GetData()) == 0 {
		FatalError("bad block", errors.New("invalid block"))
	}

	txs := blk.GetData().GetData()
	for index, data := range txs {
		logger.Info("tx", zap.Any("tx", index))
		parseTx(data)
	}
}

func parseTx(envBytes []byte) { //nolint:unused
	processedTransaction := peer.ProcessedTransaction{}
	err := proto.Unmarshal(envBytes, &processedTransaction)
	if err != nil {
		FatalError("unmarshal of transaction proposal processedTransaction failed", err)
	}
	logger.Info("processedTransaction.ValidationCode", zap.Any("ValidationCode", processedTransaction.GetValidationCode()))

	txID, err := GetOrComputeTxIDFromEnvelope(envBytes)
	if err != nil {
		FatalError("error GetOrComputeTxIDFromEnvelope", err)
	}
	logger.Info("txID", zap.Any("txID", txID))

	transaction, err := UnmarshalTransaction(envBytes)
	if err != nil {
		FatalError("error UnmarshalTransaction", err)
	}

	for i, transactionAction := range transaction.GetActions() {
		logger.Info("transactionAction ", zap.Any("i", i))
		rwSetByTransactionAction(transactionAction)
	}
}

func rwSetByTransactionAction(transactionAction *peer.TransactionAction) { //nolint:unused
	_, ca, err := GetPayloads(transactionAction)
	if err != nil {
		return
	}

	chaincodeEvent := &peer.ChaincodeEvent{}
	if err = proto.Unmarshal(ca.GetEvents(), chaincodeEvent); err != nil {
		return
	}

	logger.Info("chaincodeEvent.EventName", zap.Any("EventName", chaincodeEvent.GetEventName()))
}

func rwSetByChaincodeAction(chaincodeAction *peer.ChaincodeAction) { //nolint:unused
	if chaincodeAction.GetResponse().GetStatus() != http.StatusOK {
		return
	}

	chaincodeEvents, err := UnmarshalChaincodeEvents(chaincodeAction.GetEvents())
	if err != nil {
		FatalError("chaincode events", err)
	}

	logger.Info("chaincodeEvents.EventName", zap.Any("chaincodeEvents.EventName", chaincodeEvents.GetEventName()))

	txRWSet := &rwsetutil.TxRwSet{}
	err = txRWSet.FromProtoBytes(chaincodeAction.GetResults())
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
	if len(set.KvRwSet.GetReads()) == 0 {
		return
	}
	logger.Info("KvRwSet.Reads", zap.Any("Reads", set.KvRwSet.GetReads()))
}

func printWrites(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.GetWrites()) == 0 {
		return
	}
	logger.Info("KvRwSet.Writes", zap.Any("Writes", set.KvRwSet.GetWrites()))
}

func printRangeQueriesInfo(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.GetRangeQueriesInfo()) == 0 {
		return
	}
	logger.Info("KvRwSet.RangeQueriesInfo", zap.Any("RangeQueriesInfo", set.KvRwSet.GetRangeQueriesInfo()))
}

func printMetadataWrites(set *rwsetutil.NsRwSet) {
	if len(set.KvRwSet.GetMetadataWrites()) == 0 {
		return
	}
	logger.Info("KvRwSet.MetadataWrites", zap.Any("MetadataWrites", set.KvRwSet.GetMetadataWrites()))
}

func printCollectionName(set *rwsetutil.NsRwSet) {
	if len(set.CollHashedRwSets) == 0 {
		return
	}
	logger.Info("set.CollectionName", zap.Any("CollHashedRwSets", set.CollHashedRwSets))
}

// GetPayloads gets the underlying payload objects in a TransactionAction.
func GetPayloads(txActions *peer.TransactionAction) (*peer.ChaincodeActionPayload, *peer.ChaincodeAction, error) {
	ccPayload, err := GetChaincodeActionPayload(txActions.GetPayload())
	if err != nil {
		return nil, nil, err
	}

	pRespPayload, err := GetProposalResponsePayload(
		ccPayload.GetChaincodeProposalPayload(),
	)
	if err != nil {
		return nil, nil, err
	}

	logger.Info("pRespPayload", zap.Any("pRespPayload", pRespPayload))

	if pRespPayload.GetExtension() == nil {
		return nil, nil, errors.New("extension is nil for txActions.Payload.ChaincodeProposalPayload")
	}

	respPayload, err := GetChaincodeAction(pRespPayload.GetExtension())
	if err != nil {
		return nil, nil, err
	}

	payload, err := GetChaincodeActionPayload(respPayload.GetResults())
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
		zap.Any("payload.ChaincodeProposalPayload", payload.GetChaincodeProposalPayload()),
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
