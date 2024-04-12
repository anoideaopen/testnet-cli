package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/anoideaopen/testnet-cli/utils"
	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
)

var signMessageCmd = &cobra.Command{
	Use:   "signMessage",
	Short: "sign message by validator - for acl",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile("message.txt")
		if err != nil {
			FatalError("message.txt", err)
		}

		privateKey, publicKey, err := utils.GetPrivateKey(config.SecretKey)
		if err != nil {
			msg := fmt.Sprintf("Failed to GetPrivateKeySK %s", config.SecretKey)
			FatalError(msg, err)
		}
		signerInfo := utils.SignerInfo{}
		signerInfo.PublicKey = publicKey
		signerInfo.PrivateKey = privateKey

		result := strings.Split(string(data), "\n")
		signatureBytes, _, err := utils.SignMessage(signerInfo, result)
		if err != nil {
			FatalError("Error SignMessage", err)
		}
		signature := base58.Encode(signatureBytes)

		// save to file
		file, err := os.Create(fmt.Sprintf("signature-%s.txt", base58.Encode(signerInfo.PublicKey)))
		if err != nil {
			FatalError("Error create signature-%s.txt", err)
		}
		defer file.Close()

		_, err = file.WriteString(signature)
		if err != nil {
			FatalError("signature-%s.txt WriteString", err)
		}
	},
}
