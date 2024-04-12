package report

import (
	"errors"
	"time"

	"github.com/anoideaopen/testnet-cli/db/postgres"
	"github.com/go-pg/pg/v10"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	_ "github.com/lib/pq"
)

type PostgresReporter struct {
	db *pg.DB
}

func NewPostgresReporter(db *pg.DB) PostgresReporter {
	return PostgresReporter{
		db: db,
	}
}

func (p PostgresReporter) Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error {
	if resp == nil {
		return errors.New("response can'p be nil")
	}

	request := postgres.Request{
		TxID:               string(resp.TransactionID),
		BlockNumber:        resp.BlockNumber,
		CreatedAt:          timeStart.UTC().Unix(),
		Finished:           timeEnd.UTC().Unix(),
		Duration:           timeEnd.Sub(timeStart).Milliseconds(),
		ValidationCode:     int32(resp.TxValidationCode),
		ValidationCodeText: resp.TxValidationCode.String(),
		ChaincodeStatus:    resp.ChaincodeStatus,
	}

	return postgres.InsertRequest(p.db, request)
}

func (p PostgresReporter) Close() error {
	return p.db.Close()
}
