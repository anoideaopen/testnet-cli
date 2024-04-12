package cmd

import (
	"crypto/ed25519"
	"fmt"
	"io/ioutil"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "get public key by private key (public key -> base58.Check)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.PrivateKeyFilePath) != 0 {
			privateKey, err := GetPrivateKeyByFile(config.PrivateKeyFilePath)
			if err != nil {
				logger.Error("GetPrivateKeyByFile", zap.Error(err))
				return
			}

			pubk, _ := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
			fmt.Println(base58.Encode(pubk))
			return
		}
		secretKey := config.SecretKey
		if secretKey == "" && len(args) == 1 {
			secretKey = args[0]
		}
		publicKey, err := utils.GetPublicKey(secretKey)
		if err != nil {
			logger.Error("getPublicKey", zap.Error(err))
			return
		}
		fmt.Println(publicKey)
	},
}

func GetPrivateKeyByFile(privateKeyPath string) (privateKeyBytes []byte, err error) {
	privateKeyBytes, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	return privateKeyBytes, nil
}
