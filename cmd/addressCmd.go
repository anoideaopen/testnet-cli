package cmd

import (
	"crypto/ed25519"
	"fmt"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/anoideaopen/testnet-cli/utils"
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
		var publicKey ed25519.PublicKey
		if err == nil {
			var err error

			_, publicKey, err = utils.GetPrivateKey(secretKey)
			if err != nil {
				logger.Error("GetPrivateKey", zap.Error(err))
				return
			}
		} else {
			publicKey = base58.Decode(secretKey)
		}

		addr, err := utils.GetAddressByPublicKey(publicKey)
		if err != nil {
			logger.Error("GetAddressByPublicKey", zap.Error(err))
			return
		}
		fmt.Println(addr)
	},
}
