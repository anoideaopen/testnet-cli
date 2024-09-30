package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var sendRequestCmd = &cobra.Command{ //nolint:unused
	Use:   "sendRequest",
	Short: "send to HLF generated message with validator's signatures",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID, methodName, methodArgs := handlerArgs(args)

		logger.Debug(channelID)
		logger.Debug(methodName)
		logger.Debug(fmt.Sprintf("%v\n", methodArgs))

		signatureFilePaths := methodArgs[0]

		data, err := os.ReadFile("message.txt")
		if err != nil {
			return
		}

		logger.Debug("read data from message.txt")
		var result []string
		for _, s := range strings.Split(string(data), "\n") {
			result = append(result, s)
			logger.Debug(s)
		}

		var signatures []string
		for _, signatureFilePath := range strings.Split(signatureFilePaths, ",") {
			filename := fmt.Sprintf("signature-%s.txt", signatureFilePath)
			data, err := os.ReadFile(filename)
			logger.Debug(fmt.Sprintf("filename %s\n", filename))
			if err != nil {
				return
			}
			logger.Debug("data: " + string(data))
			signatures = append(signatures, string(data))
		}

		messageWithSigArg := append(result[1:], signatures...)

		logger.Debug(
			"Sign result",
			zap.Strings("messageWithSig", messageWithSigArg),
		)

		for i, s := range messageWithSigArg {
			logger.Info(fmt.Sprintf("message %d\n", i))
			logger.Info(s)
		}

		logger.Info(fmt.Sprintf("channelID %s\n", channelID))
		logger.Info(fmt.Sprintf("methodName %s\n", methodName))

		logger.Debug("send request to HLF...")
		response, err := HlfClient.Invoke(config.WaitBatch, channelID, config.ChaincodeName, methodName, messageWithSigArg)
		if err != nil {
			FatalError("Invoke", err)
		}

		fmt.Println("response.Responses:")
		fmt.Println(response.Responses)
		fmt.Println("response.TransactionID:")
		fmt.Println(response.TransactionID)
	},
}
