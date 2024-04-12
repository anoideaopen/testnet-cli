package cmd

import (
	"fmt"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var privkeyCmd = &cobra.Command{
	Use:   "privkey",
	Short: "generate private key (private key -> base58.CheckEncode)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		isHex := false
		if len(args) == 1 && args[0] == "hex" {
			isHex = true
		}

		if len(config.PrivateKeyFilePath) != 0 {
			privateKey, err := GetPrivateKeyByFile(config.PrivateKeyFilePath)
			if err != nil {
				logger.Error("GetPrivateKeyByFile", zap.Error(err))
				return
			}

			if isHex {
				privateKeyHex := utils.ConvertPrivateKeyToHex(privateKey)
				fmt.Println(privateKeyHex)
			} else {
				privateKeyBase58Check := utils.ConvertPrivateKeyToBase58Check(privateKey)
				fmt.Println(privateKeyBase58Check)
			}
			return
		}

		_, privateKey, err := utils.GeneratePrivateAndPublicKey()
		if err != nil {
			logger.Error("generatePrivateKey", zap.Error(err))
			return
		}
		if isHex {
			privateKeyHex := utils.ConvertPrivateKeyToHex(privateKey)
			fmt.Println(privateKeyHex)
		} else {
			privateKeyBase58Check := utils.ConvertPrivateKeyToBase58Check(privateKey)
			fmt.Println(privateKeyBase58Check)
		}
	},
}
