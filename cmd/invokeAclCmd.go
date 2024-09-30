package cmd

import (
	"fmt"
	"strings"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// const noBatch = true

var invokeACLCmd = &cobra.Command{ //nolint:unused
	Use:   "invokeAcl channelID methodName [optional method arguments]",
	Short: "invoke acl version with signature in hex - v0.8.1-0.0.2 and earlier",
	Args:  cobra.MinimumNArgs(2), //nolint:gomnd
	Run: func(cmd *cobra.Command, args []string) {
		initHlfClient()

		channelID, methodName, methodArgs := handlerArgs(args)

		logger.Debug(channelID)
		logger.Debug(methodName)
		fmt.Printf("%v\n", methodArgs)

		address := methodArgs[0]
		reason := methodArgs[1]
		reasonID := methodArgs[2]
		newPkey := methodArgs[3]

		logger.Debug("methodArgs")
		for i, arg := range methodArgs {
			fmt.Printf("[%d]\n", i)
			fmt.Printf("    - '%v'\n", arg)
		}

		var validators []utils.SignerInfo
		validatorsKey := strings.Split(config.SecretKey, ",")
		for _, secretKey := range validatorsKey {
			logger.Info("secretKey", zap.String("secretKey", secretKey))
			privateKey, publicKey, err := utils.GetPrivateKey(secretKey)
			if err != nil {
				msg := "Failed to GetPrivateKeySK " + secretKey
				FatalError(msg, err)
			}
			signerInfo := utils.SignerInfo{}
			signerInfo.PublicKey = publicKey
			signerInfo.PrivateKey = privateKey
			validators = append(validators, signerInfo)
		}

		signedMessageArg, _, err := utils.SignACL(validators, methodName, address, reason, reasonID, newPkey)
		logger.Debug("--- signedMessage")
		for i, arg := range signedMessageArg {
			fmt.Printf("%d\n", i)
			fmt.Printf("%v\n", arg)
		}
		if err != nil {
			FatalError("err signedMessage", err)
		}

		response, err := HlfClient.Invoke(false, channelID, config.ChaincodeName, methodName, signedMessageArg)
		if err != nil {
			FatalError("Invoke", err)
		}

		fmt.Println("response.Responses:")
		fmt.Println(response.Responses)
		fmt.Println("response.TransactionID:")
		fmt.Println(response.TransactionID)
	},
}
