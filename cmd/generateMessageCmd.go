package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/spf13/cobra"
)

var generateMessageCmd = &cobra.Command{
	Use:   "generateMessage",
	Short: "generate message for validators - for acl",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("secretKey")
		logger.Debug(config.SecretKey)

		channelID, methodName, methodArgs := handlerArgs(args)

		logger.Debug(channelID)
		logger.Debug(methodName)
		logger.Debug(fmt.Sprintf("%v\n", methodArgs))

		logger.Debug("methodArgs")
		for i, arg := range methodArgs {
			logger.Debug(fmt.Sprintf("[%d]\n", i))
			logger.Debug(fmt.Sprintf("    - '%v'\n", arg))
		}

		validatorPublicKeys := methodArgs[0]

		message := utils.GenerateMessage(strings.Split(validatorPublicKeys, ","), channelID, config.ChaincodeName, methodName, methodArgs[1:])
		file, err := os.Create("message.txt")
		if err != nil {
			return
		}
		defer file.Close()

		_, err = file.WriteString(message)
		if err != nil {
			panic(err)
		}
	},
}
