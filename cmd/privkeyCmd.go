package cmd

import (
	"fmt"

	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/service"
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

		keyType := proto.KeyType(config.KeyType)

		if len(config.PrivateKeyFilePath) != 0 {
			privateKey, err := GetPrivateKeyByFile(config.PrivateKeyFilePath)
			if err != nil {
				logger.Error("GetPrivateKeyByFile", zap.Error(err))
				return
			}

			if isHex {
				privateKeyHex := service.BytesToHex(privateKey)
				fmt.Println(privateKeyHex)
			} else {
				privateKeyBase58Check := service.ConvertPrivateKeyToBase58CheckFromBytes(privateKey)
				fmt.Println(privateKeyBase58Check)
			}
			return
		}

		keys, err := service.GeneratePrivateAndPublicKey(keyType)
		if err != nil {
			logger.Error("generatePrivateKey", zap.Error(err))
			return
		}
		if isHex {
			privateKeyHex, err := service.ConvertPrivateKeyToHex(keys)
			if err != nil {
				logger.Error("ConvertPrivateKeyToHex", zap.Error(err))
			}
			fmt.Println(privateKeyHex)
		} else {
			privateKeyBase58Check, err := service.ConvertPrivateKeyToBase58Check(keys)
			if err != nil {
				logger.Error("ConvertPrivateKeyToBase58Check", zap.Error(err))
			}
			fmt.Println(privateKeyBase58Check)
		}
	},
}
