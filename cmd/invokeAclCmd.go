package cmd

import (
	"fmt"
	"strings"

	"github.com/anoideaopen/foundation/keys"
	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// const noBatch = true

var invokeACLCmd = &cobra.Command{
	Use:   "invokeAcl channelID methodName [optional method arguments]",
	Short: "invoke acl version with signature in hex - v0.8.1-0.0.2 and earlier",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID, methodName, methodArgs := handlerArgs(args)

		logger.Debug("channelID", zap.String("channelID", channelID))
		logger.Debug("methodName", zap.String("methodName", methodName))
		logger.Debug("methodArgs", zap.Any("methodArgs", methodArgs))

		for i, arg := range methodArgs {
			fmt.Printf("[%d] '%v'\n", i, arg)
		}

		var reqArgs []string
		if config.SecretKey != "" {
			validatorsKey := strings.Split(config.SecretKey, ",")
			validators := make([]*keys.Keys, 0, len(validatorsKey))
			keyType := proto.KeyType(config.KeyType)

			for _, secretKey := range validatorsKey {
				logger.Info("secretKey", zap.String("secretKey", secretKey))

				k, err := service.GetKeys(secretKey, keyType)
				if err != nil {
					FatalError("Failed to GetPrivateKey "+secretKey, err)
				}

				validators = append(validators, k)
			}

			var err error
			reqArgs, _, err = service.SignACL(validators, methodName, methodArgs)
			if err != nil {
				FatalError("Failed to sign ACL", err)
			}

			fmt.Println("Signed message arguments:")
		} else {
			reqArgs = methodArgs
			fmt.Println("Unsigned message arguments:")
		}

		for i, arg := range reqArgs {
			fmt.Printf("[%d] %v\n", i, arg)
		}

		response, err := HlfClient.Invoke(false, channelID, config.ChaincodeName, methodName, reqArgs)
		if err != nil {
			FatalError("Invoke", err)
		}

		fmt.Println("response.Responses:")
		fmt.Println(response.Responses)
		fmt.Println("response.TransactionID:")
		fmt.Println(response.TransactionID)
	},
}
