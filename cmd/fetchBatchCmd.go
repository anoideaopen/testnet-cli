package cmd

import (
	"context"
	"sync"

	"github.com/anoideaopen/testnet-cli/db/postgres"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/observer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var fetchBatchCmd = &cobra.Command{ //nolint:unused
	Use:   "fetchBatch",
	Short: "fetchBatch",
	Args:  cobra.MinimumNArgs(0), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		db, err := postgres.NewPostgresDB(
			config.Postgres.PostgresHost,
			config.Postgres.PostgresPort,
			config.Postgres.PostgresUser,
			config.Postgres.PostgresPassword,
			config.Postgres.PostgresDBName,
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				panic(err)
			}
		}()
		observer := observer.NewObserver(
			config.Observer.ObserverUsername,
			config.Observer.ObserverPassword,
			config.Observer.ObserverURL,
			config.Observer.ObserverVersion,
		)
		requests := postgres.GetRequests(db)

		// Create a rate limiter with a limit of 100 requests per second
		limiter := rate.NewLimiter(rate.Limit(config.RequestsPerSecond), 1)

		numberRequest := len(requests)

		// WaitGroup to wait for all goroutines to finish
		var wg sync.WaitGroup
		wg.Add(numberRequest)

		ctx := context.Background()
		for i := range numberRequest {
			// Wait for the rate limiter to allow the next request
			err := limiter.WaitN(ctx, 1)
			if err != nil {
				logger.GetLogger().Error("batch", zap.Error(err))
				return
			}

			go func(tx string) {
				defer wg.Done()
				batch, err := observer.GetBatch(tx)
				if err != nil {
					logger.GetLogger().Error("batch", zap.Error(err))
					return
				} else {
					if err := postgres.InsertBatch(db, batch); err != nil {
						panic(err)
					}
				}
			}(requests[i].TxID)
		}
		// Wait for all goroutines to finish
		wg.Wait()
	},
}
