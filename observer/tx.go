package observer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anoideaopen/testnet-cli/db/postgres"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/go-pg/pg/v10"
	"go.uber.org/zap"
)

type Tx struct {
	Txs []struct {
		Header struct {
			TxID      string `json:"tx_id"`
			Timestamp any    `json:"timestamp"`
		} `json:"header"`
		Metadata struct {
			Error          string `json:"error"`
			ValidationCode int    `json:"validation_code"`
			Number         int    `json:"number"`
			BatchID        string `json:"batch_id"`
		} `json:"metadata"`
	} `json:"txs"`
}

const (
	layout = "2006-01-02T15:04:05.999999"
)

type Observer struct {
	username string
	password string
	url      string
	version  string
}

func NewObserver(
	username string,
	password string,
	url string,
	version string,
) Observer {
	return Observer{
		username: username,
		password: password,
		url:      url,
		version:  version,
	}
}

func (o Observer) TxAPI(txID string) (Tx, error) {
	if o.url == "" {
		err := errors.New("observer url empty")
		logger.Error("", zap.Error(err))
		return Tx{}, err
	}

	if o.version == "" {
		err := errors.New("observer version empty")
		logger.Error("", zap.Error(err))
		return Tx{}, err
	}

	if txID == "" {
		err := errors.New("observer txID empty")
		logger.Error("", zap.Error(err))
		return Tx{}, err
	}

	client := &http.Client{
		Transport: &http.Transport{},
	}

	req, err := http.NewRequest("GET", o.url+"/"+o.version+"/tx/"+txID, nil)
	if err != nil {
		fmt.Println("Failed prepare request", err)
		return Tx{}, err
	}

	if o.username != "" {
		req.SetBasicAuth(o.username, o.password)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed send request:", err)
		return Tx{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return Tx{}, err
	}

	var tx Tx
	err = json.Unmarshal(body, &tx)
	if err != nil {
		fmt.Println("Can't unmarshal request", err)
		return Tx{}, err
	}
	return tx, nil
}

func (o Observer) GetBatch(requestTxID string) (postgres.Batch, error) {
	requestTx, err := o.TxAPI(requestTxID)
	if err != nil {
		logger.Error("TxAPI", zap.Error(err))
		return postgres.Batch{}, err
	}
	if len(requestTx.Txs) == 0 {
		err := errors.New("not found txs")
		logger.Error("TxAPI response txs not found", zap.Error(err))
		return postgres.Batch{}, err
	}

	batchID := requestTx.Txs[0].Metadata.BatchID
	if batchID == "" {
		err := fmt.Errorf("can't find batchID for request %s", requestTxID)
		logger.Debug("GetBatch", zap.Error(err))
		return postgres.Batch{}, err
	}

	firstTimestamp := requestTx.Txs[0].Header.Timestamp.(string)

	batchTx, err := o.TxAPI(batchID)
	if err != nil {
		logger.Error("TxAPI", zap.Error(err))
		return postgres.Batch{}, err
	}

	blockNumber := batchTx.Txs[0].Metadata.Number
	batchTimestamp := batchTx.Txs[0].Header.Timestamp.(string)
	batchValidationCode := batchTx.Txs[0].Metadata.ValidationCode
	batchErrMsg := batchTx.Txs[0].Metadata.Error

	firstTime, err := time.Parse(layout, firstTimestamp)
	if err != nil {
		logger.Error("Parse firstTimestamp", zap.Error(err))
		return postgres.Batch{}, err
	}

	batchTime, err := time.Parse(layout, batchTimestamp)
	if err != nil {
		logger.Error("Parse batchTimestamp", zap.Error(err))
		return postgres.Batch{}, err
	}

	return postgres.Batch{
		TxID:                batchID,
		BlockNumber:         blockNumber,
		CreatedAt:           batchTime.UTC().Unix(),
		BatchValidationCode: batchValidationCode,
		RequestTxID:         requestTxID,
		RequestCreatedAt:    firstTime.UTC().Unix(),
		BatchErrorMsg:       batchErrMsg,
	}, nil
}

func (o Observer) FetchBatches(db *pg.DB, requests []postgres.Request) {
	for _, request := range requests {
		batch, err := o.GetBatch(request.TxID)
		if err != nil {
			logger.GetLogger().Error("get batch from observer", zap.Error(err))
			continue
		} else {
			if err := postgres.InsertBatch(db, batch); err != nil {
				panic(err)
			}
		}
	}
}
