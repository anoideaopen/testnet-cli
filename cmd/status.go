package cmd

import (
	"fmt"
	"time"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/observer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "status",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		observer := observer.NewObserver(
			config.Observer.ObserverUsername,
			config.Observer.ObserverPassword,
			config.Observer.ObserverURL,
			config.Observer.ObserverVersion,
		)
		if len(args) != 1 || args[0] == "" {
			logger.GetLogger().Error("tx id required")
			return
		}
		txID := args[0]
		batch, err := observer.GetBatch(txID)
		if err != nil {
			logger.GetLogger().Error("get batch from observer", zap.Error(err))
		} else {
			fmt.Println("\n\n\n-------- Batch tx found in observer:")
			fmt.Println("Request:")
			fmt.Printf("TxID: %s\n", batch.RequestTxID)
			fmt.Printf("CreatedAt: %s\n\n", time.Unix(time.Now().Unix(), batch.RequestCreatedAt).Format(time.RFC3339Nano))

			fmt.Println("Batch:")
			fmt.Printf("TxID: %s\n", batch.TxID)
			fmt.Printf("BlockNumber: %d\n", batch.BlockNumber)
			fmt.Printf("CreatedAt: %s\n", time.Unix(time.Now().Unix(), batch.CreatedAt).Format(time.RFC3339Nano))
			fmt.Printf("BatchErrorMsg: %s\n", batch.BatchErrorMsg)
			fmt.Printf("BatchValidationCode: %d\n", batch.BatchValidationCode)
		}
	},
}
