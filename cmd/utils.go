package cmd

import (
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/service"
	"go.uber.org/zap"
)

func handlerArgs(args []string) (string, string, []string) {
	channelID := args[0]
	methodName := args[1]
	methodArgs := args[2:]

	if len(config.ChaincodeName) == 0 {
		config.ChaincodeName = channelID
	}

	return channelID, methodName, methodArgs
}

func FatalError(errorMessage string, err error) {
	if err != nil {
		logger.Error(errorMessage, zap.Error(err))
	} else {
		logger.Error(errorMessage)
	}

	panic(err)
}

var HlfClient *service.HLFClient

func initHlfClient() {
	var err error
	HlfClient, err = service.NewHLFClient(config.Connection, config.User, config.Organization, nil)
	if err != nil {
		FatalError("Failed to create new channel client", err)
	}
}
