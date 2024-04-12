package report

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type ResponseReporter struct {
}

func (t ResponseReporter) Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error {
	if resp != nil && resp.TxValidationCode == 0 {
		isQuery := resp.BlockNumber == 0
		if !isQuery {
			fmt.Println("TransactionID:")
			fmt.Println(resp.TransactionID)
			fmt.Println("TxValidationCode:")
			fmt.Println(resp.TxValidationCode)
			fmt.Println("BlockNumber:")
			fmt.Println(resp.BlockNumber)
			payload := string(resp.Payload)
			fmt.Println(payload)
		}
		payload := string(resp.Payload)
		fmt.Println(payload)
	} else {
		return fmt.Errorf("error resp: %v", resp)
	}
	return nil
}

func (t ResponseReporter) Close() error {
	return nil
}
