package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto" //nolint:staticcheck
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/spf13/cobra"
)

var getTxIDFromBlockCmd = &cobra.Command{
	Use:   "getTxIDFromBlock filePath",
	Short: "compute requestMaxBytes by block file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		if len(filePath) == 0 {
			FatalError("block file path is empty", nil)
		}

		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			FatalError("failed ReadFile", err)
		}

		block := &cb.Block{}
		err = proto.Unmarshal(bytes, block)
		if err != nil {
			FatalError("failed Unmarshal", err)
		}

		for _, envBytes := range block.Data.Data {
			env, err := protoutil.GetEnvelopeFromBlock(envBytes)
			if err != nil {
				FatalError("GetEnvelopeFromBlock", err)
			}

			payload, err := protoutil.UnmarshalPayload(env.Payload)
			if err != nil {
				FatalError("UnmarshalPayload", err)
			}

			chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
			if err != nil {
				FatalError("UnmarshalChannelHeader", err)
			}

			fmt.Println(chdr.TxId)
			fmt.Println(payload.Data)
		}
	}}
