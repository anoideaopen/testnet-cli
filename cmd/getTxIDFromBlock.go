package cmd

import (
	"fmt"
	"os"

	cb "github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
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

		bytes, err := os.ReadFile(filePath)
		if err != nil {
			FatalError("failed ReadFile", err)
		}

		block := &cb.Block{}
		err = proto.Unmarshal(bytes, block)
		if err != nil {
			FatalError("failed Unmarshal", err)
		}

		for _, envBytes := range block.GetData().GetData() {
			env, err := GetEnvelopeFromBlock(envBytes)
			if err != nil {
				FatalError("GetEnvelopeFromBlock", err)
			}

			payload, err := UnmarshalPayload(env.GetPayload())
			if err != nil {
				FatalError("UnmarshalPayload", err)
			}

			chdr, err := UnmarshalChannelHeader(payload.GetHeader().GetChannelHeader())
			if err != nil {
				FatalError("UnmarshalChannelHeader", err)
			}

			fmt.Println(chdr.GetTxId())
			fmt.Println(payload.String())
		}
	},
}
