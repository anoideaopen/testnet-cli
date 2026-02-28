package cmd

import (
	"encoding/hex"
	"time"

	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/report"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/spf13/cobra"
	pb "google.golang.org/protobuf/proto"
)

var batchExecuteCmd = &cobra.Command{
	Use:   "batchExecute channelID tx1 tx2 tx3",
	Short: "execute method batchExecute in chaincode",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		// type of report response for example print in log or store metric in InfluxDB
		reporter := report.GetReporterByID(config)

		channelID := args[0]
		txIDs := args[1:]

		batch := proto.Batch{}
		for _, txID := range txIDs {
			binaryTx, err := hex.DecodeString(txID)
			if err != nil {
				panic(err)
			}
			batch.TxIDs = append(batch.TxIDs, binaryTx)
		}
		b, _ := pb.Marshal(&batch)
		reqArgs := []string{string(b)}

		requestOptions := prepareRequestOptions()
		timeout := channel.WithTimeout(fab.Execute, 180*time.Second)
		requestOptions = append(requestOptions, timeout)
		start := time.Now()
		response, err := HlfClient.Query(channelID, channelID, "batchExecute", reqArgs, requestOptions...)
		end := time.Now()
		if err != nil {
			panic(err)
		}

		err = reporter.Report(response, start, end)
		if err != nil {
			FatalError("Query", err)
			return
		}
	},
}
