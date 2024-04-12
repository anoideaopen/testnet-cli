package report

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type TableReporter struct {
}

func (t TableReporter) Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error {
	if resp.TxValidationCode == 0 {
		fmt.Printf("| tx | %s | block | %d | start | %s | end | %s | dur | %f |\n",
			string(resp.TransactionID),
			resp.BlockNumber,
			timeStart.UTC().Format("2006-01-02 15:04:05.000000"),
			timeEnd.UTC().Format("2006-01-02 15:04:05.000000"),
			timeEnd.Sub(timeStart).Seconds(),
		)
	}
	return nil
}

func (t TableReporter) Close() error {
	return nil
}
