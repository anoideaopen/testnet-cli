package report

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type EndorserReporter struct{}

func (t EndorserReporter) Report(resp *channel.Response, timeStart time.Time, timeEnd time.Time) error {
	fmt.Println("-------- Request transaction info:")
	fmt.Println("TransactionID:")
	fmt.Println(resp.TransactionID)
	fmt.Println("TxValidationCode:")
	fmt.Println(resp.TxValidationCode)
	fmt.Println("BlockNumber:")
	fmt.Println(resp.BlockNumber)
	fmt.Println("Endorser responses:")
	for i, r := range resp.Responses {
		fmt.Printf("\n\n- [%d] Response: \n", i)
		fmt.Printf("Endorser %s\n", r.Endorser)
		fmt.Printf("Status: %d\n", r.Response.GetStatus())
		// fmt.Printf("Message: %s\n", r.Response.Message)
		// fmt.Printf("Payload: %s\n", string(r.Response.Payload))
	}
	return nil
}

func (t EndorserReporter) Close() error {
	return nil
}
