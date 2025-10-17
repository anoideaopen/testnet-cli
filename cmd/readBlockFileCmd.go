package cmd

// import (
//	"errors"
//	"fmt"
//	"os"
//
//	"github.com/golang/protobuf/proto"
//	cb "github.com/hyperledger/fabric-protos-go-apiv2/common"
//	"github.com/spf13/cobra"
// )
//
// var (
//	readBlockFileCmd = &cobra.Command{
//		Use:   "readBlockFile channelID blockID peer",
//		Short: "readBlock by blockID by file",
//		Args:  cobra.ExactArgs(3),
//		Run: func(cmd *cobra.Command, args []string) {
//			logger.GetLogger()
//
//			channelID := args[0]
//			if len(channelID) == 0 {
//				FatalError("channelID is empty", nil)
//			}
//
//			blockID := args[1]
//			if len(blockID) == 0 {
//				FatalError("blockID is empty", nil)
//			}
//
//			peer := args[2]
//			if len(peer) == 0 {
//				FatalError("peer is empty", nil)
//			}
//
//			block := &cb.Block{}
//
//			file := channelID + "_" + blockID + ".block"
//			bytes, err := os.ReadFile(file)
//			if err != nil {
//				FatalError("failed ReadFile", err)
//			}
//
//			if len(bytes) == 0 {
//				initHlfClient()
//				fmt.Println("========================")
//				fmt.Println("GET BLOCK")
//				fmt.Println("========================")
//				block, err = getBlock(channelID, blockID, peer)
//				if err != nil {
//					FatalError("failed to get block", err)
//				}
//				if block == nil {
//					FatalError("block is nil", errors.New(fmt.Sprintf("failed to get block. channelID %s, blockID %s peer %s", channelID, blockID, peer)))
//				}
//				err = saveBlock(channelID, blockID, block)
//				if err != nil {
//					FatalError("failed to save block", err)
//				}
//			} else {
//				logger.Info("================================")
//				logger.Info("READ file with block")
//				logger.Info("================================")
//				err = proto.Unmarshal(bytes, block)
//				if err != nil {
//					FatalError("failed Unmarshal", err)
//				}
//			}
//
//			printBlock2(channelID, block)
//		}}
// )
//
// //
// //func getBlock(channelID string, blockID string, peer string) (*cb.Block, error) {
// //	block, err := HlfClient.QueryBlock(channelID, blockID, peer)
// //	if err != nil {
// //		FatalError("Failed to QueryBlock", err)
// //		return nil, err
// //	}
// //
// //	return block, nil
// //}
// //
// //func saveBlock(channelID string, blockID string, block *cb.Block) error {
// //	logger.Info("================================")
// //	logger.Info("Start save block")
// //	logger.Info("================================")
// //	b, err := proto.Marshal(block)
// //	if err != nil {
// //		return err
// //	}
// //
// //	file := channelID + "_" + blockID + ".block"
// //	if err = os.WriteFile(file, b, 0644); err != nil {
// //		return err
// //	}
// //
// //	logger.Info("================================")
// //	logger.Info("End save block")
// //	logger.Info("================================")
// //	return nil
// //}
// //
// type TxValidationFlags []uint8
//
// func printBlock2(channelID string, block *cb.Block) {
//	logger.Info("================================")
//	logger.Info("Print block")
//	logger.Info("================================")
//
//	codes := TxValidationFlags(block.Metadata.Metadata[cb.BlockMetadataIndex_TRANSACTIONS_FILTER])
//
//	for txIndex, txData := range block.Data.Data {
//		logger.Info("\n=========================")
//		logger.Info(fmt.Sprintf("%v", txIndex))
//		logger.Info("\n=========================")
//
//		envelope := &cb.Envelope{}
//		if err := proto.Unmarshal(txData, envelope); err != nil {
//			panic(err)
//		}
//
//		tx, err := hlf.ExtractTx(channelID, block.Header.Number, envelope.Payload, int32(codes[txIndex]))
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(tx)
//
//		//var env *cb.Envelope
//		//var payload *cb.Payload
//		//var err error
//		//if env, err = GetEnvelopeFromBlock(envBytes); err == nil {
//		//	if payload, err = UnmarshalPayload(env.Payload); err == nil {
//		//		shdr, err := UnmarshalSignatureHeader(payload.Header.SignatureHeader)
//		//		if err != nil {
//		//			FatalError("UnmarshalSignatureHeader", err)
//		//		}
//		//		logger.Info("shdr.Creator")
//		//		logger.Info(string(shdr.Creator))
//		//		serializedIdentity, err := UnmarshalSerializedIdentity(shdr.Creator)
//		//		if err != nil {
//		//			FatalError("UnmarshalSerializedIdentity", err)
//		//		}
//		//
//		//		logger.Info("serializedIdentity.Mspid")
//		//		logger.Info(serializedIdentity.Mspid)
//		//
//		//		logger.Info("serializedIdentity.IdBytes")
//		//		logger.Info(string(serializedIdentity.IdBytes))
//		//
//		//		logger.Info("shdr.Nonce")
//		//		logger.Info(string(shdr.Nonce))
//		//
//		//		chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
//		//		if err != nil {
//		//			FatalError("UnmarshalChannelHeader", err)
//		//		}
//		//		logger.Info("chdr.TxId")
//		//		logger.Info(chdr.TxId)
//		//	}
//		//}
//	}
// }
