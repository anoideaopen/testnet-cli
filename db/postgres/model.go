package postgres

import (
	"fmt"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"go.uber.org/zap"
)

type Batch struct {
	TxID                string `json:"tx_id,omitempty"`
	BlockNumber         int    `json:"block_number,omitempty"`
	CreatedAt           int64  `json:"created_at,omitempty"`
	BatchValidationCode int    `json:"batch_validation_code,omitempty"`
	BatchErrorMsg       string `json:"batch_error_msg,omitempty"`
	RequestTxID         string `json:"request_tx_id,omitempty"`
	RequestCreatedAt    int64  `json:"request_created_at,omitempty"`
}

func (b Batch) String() string {
	return fmt.Sprintf("Batch<%s %d %d %d %s %d>",
		b.TxID, b.BlockNumber, b.CreatedAt,
		b.BatchValidationCode, b.RequestTxID, b.RequestCreatedAt)
}

type Request struct {
	TxID               string `json:"tx_id,omitempty"`
	BlockNumber        uint64 `json:"block_number,omitempty"`
	CreatedAt          int64  `json:"created_at,omitempty"`
	Finished           int64  `json:"finished,omitempty"`
	Duration           int64  `json:"duration,omitempty"`
	ValidationCode     int32  `json:"validation_code,omitempty"`
	ValidationCodeText string `json:"validation_code_text,omitempty"`
	ChaincodeStatus    int32  `json:"chaincode_status,omitempty"`
}

func (b Request) String() string {
	return fmt.Sprintf("Request<%s %d %d %d %d %d %s %d>",
		b.TxID, b.BlockNumber, b.CreatedAt,
		b.Finished, b.Duration, b.ValidationCode,
		b.ValidationCodeText, b.ChaincodeStatus)
}

// CreateSchema creates database schema for Batch and Request models.
func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*Batch)(nil),
		(*Request)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRequests(db *pg.DB) []Request {
	requestTxID := make([]string, 0)
	err := db.Model(&Batch{}).Column("request_tx_id").Select(&requestTxID)
	if err != nil {
		panic(err)
	}
	requests := make([]Request, 0)
	err = db.Model(&requests).Where("tx_id NOT IN (?)", requestTxID).Select()
	if err != nil {
		panic(err)
	}

	return requests
}

func InsertBatch(db *pg.DB, b Batch) error {
	if b.TxID == "" {
		return nil
	}
	if b.RequestTxID == "" {
		return nil
	}

	_, err := db.Model(&b).Insert()
	if err != nil {
		logger.GetLogger().Error("insert batch", zap.Error(err))
		return err
	}

	return nil
}

func InsertRequest(db *pg.DB, r Request) error {
	if r.TxID == "" {
		return nil
	}
	if r.ValidationCodeText == "" {
		return nil
	}
	_, err := db.Model(&r).Insert()
	if err != nil {
		logger.GetLogger().Error("insert request", zap.Error(err))
		return err
	}

	return nil
}
