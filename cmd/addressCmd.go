package cmd

import (
	"fmt"

	"github.com/anoideaopen/testnet-cli/service"

	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var addressCmd = &cobra.Command{
	Use:   "address",
	Short: "get address by private key (public key -> sha3.Sum256 -> base58.CheckEncode)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		secretKey := config.SecretKey
		if secretKey == "" && len(args) == 1 {
			secretKey = args[0]
		}

		_, _, err := base58.CheckDecode(secretKey)
		var publicKey string
		if err == nil {
			var err error

			keyType := proto.KeyType(config.KeyType)
			publicKey, err = service.GetPublicKey(secretKey, keyType)
			if err != nil {
				logger.Error("GetPrivateKey", zap.Error(err))
				return
			}
		} else {
			publicKey = secretKey
		}

		addr, err := service.GetAddressByPublicKey(publicKey)
		if err != nil {
			logger.Error("GetAddressByPublicKey", zap.Error(err))
			return
		}
		fmt.Println(addr)
	},
}
