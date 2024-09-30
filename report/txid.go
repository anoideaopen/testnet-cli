package report

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type TxIDReporter struct{}

func (t TxIDReporter) Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error {
	if resp != nil {
		fmt.Println(resp.TransactionID)
	}
	return nil
}

func (t TxIDReporter) Close() error {
	return nil
}
