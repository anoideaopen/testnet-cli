package cmd

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/anoideaopen/testnet-cli/db/postgres"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/observer"
	"github.com/anoideaopen/testnet-cli/report"
	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var queryCmd = &cobra.Command{
	Use:   "query channelID methodName [ optional method arguments]",
	Short: "query chaincode (with or without signatures)",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()
		channelID, methodName, methodArgs := handlerArgs(args)

		// type of report response for example print in log or store metric in InfluxDB
		reporter := report.GetReporterByID(config)
		defer func() {
			err := reporter.Close()
			if err != nil {
				panic(err)
			}
		}()

		// Create a rate limiter with a limit of 100 requests per second
		limiter := rate.NewLimiter(rate.Limit(config.RequestsPerSecond), 1)

		requestOptions := prepareRequestOptions()

		var err error
		var privateKey ed25519.PrivateKey
		var publicKey ed25519.PublicKey
		if config.SecretKey != "" {
			privateKey, publicKey, err = utils.GetPrivateKey(config.SecretKey)
			if err != nil {
				logger.Error("failed getPrivateKey", zap.Error(err))
				return
			}

			if len(privateKey) == 0 {
				logger.Error("privateKey can't be empty")
				return
			}
			if len(publicKey) == 0 {
				logger.Error("publicKey can't be empty")
				return
			}
		}

		sumTime := float64(0)
		reqCount := 0

		// WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup
		numberRequest := config.NumberRequest
		if numberRequest == 0 {
			numberRequest = math.MaxInt32
			fmt.Println(numberRequest)
		}

		wg.Add(numberRequest)

		ctx := context.Background()
		for i := 0; i < numberRequest; i++ {
			// Wait for the rate limiter to allow the next request
			err := limiter.WaitN(ctx, 1)
			if err != nil {
				panic(err)
			}

			go func() {
				defer wg.Done()

				var reqArgs []string
				if config.SecretKey != "" {
					reqArgs, err = HlfClient.SignArgs(channelID, config.ChaincodeName, methodName, methodArgs, privateKey, publicKey)
					if err != nil {
						logger.Error("failed signArgs", zap.Error(err))
						return
					}
				} else {
					reqArgs = methodArgs
				}

				start := time.Now()
				resp, err := HlfClient.Query(channelID, config.ChaincodeName, methodName, reqArgs, requestOptions...)
				if err != nil {
					fmt.Printf("Invoke error: %v\n", err)
					return
				}
				end := time.Now()

				sumTime += end.Sub(start).Seconds()
				reqCount++
				err = reporter.Report(resp, start, end)
				if err != nil {
					panic(err)
				}

				if config.WaitBatch && config.Observer.ObserverURL != "" {
					fmt.Println("\n\n\n-------- Wait batch transaction:")
					txID := string(resp.TransactionID)
					observer := observer.NewObserver(
						config.Observer.ObserverUsername,
						config.Observer.ObserverPassword,
						config.Observer.ObserverURL,
						config.Observer.ObserverVersion,
					)
					var batch postgres.Batch
					//err := retryFunc(60, 2*time.Second, func() (err error) {
					batch, err = observer.GetBatch(txID)
					//return err
					//})
					if err != nil {
						logger.GetLogger().Error("get batch from observer", zap.Error(err))
						return
					}
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
			}()
		}
		// Wait for all goroutines to finish
		wg.Wait()

		fmt.Println("request count")
		fmt.Printf("%d", reqCount)
		fmt.Println("\navg time (sec)")
		fmt.Printf("%f\n", sumTime/float64(reqCount))
	},
}
