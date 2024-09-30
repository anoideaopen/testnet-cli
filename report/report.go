package report

import (
	"time"

	"github.com/anoideaopen/testnet-cli/config"
	"github.com/anoideaopen/testnet-cli/db/postgres"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type Reporter interface {
	Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error
	Close() error
}

func GetReporterByID(config config.ApplicationConfig) Reporter {
	switch config.ResponseType {
	case "postgres":
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
		return NewPostgresReporter(db)
	case "endorsers":
		return EndorserReporter{}
	case "table":
		return TableReporter{}
	case "tx":
		return TxIDReporter{}
	default:
		return ResponseReporter{}
	}
}
