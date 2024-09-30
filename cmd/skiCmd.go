package cmd

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var skiCmd = &cobra.Command{
	Use:   "ski pathToPrivateKey",
	Short: "get ski by private key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathToPrivateKey := args[0]
		if len(pathToPrivateKey) == 0 {
			FatalError("path to private key is empty", nil)
		}

		readSKI(pathToPrivateKey)
	},
}

func readSKI(pathToPrivateKey string) {
	privateKeyFile, err := os.ReadFile(pathToPrivateKey)
	if err != nil {
		FatalError("read private key file", err)
	}
	privateKey, err := pemToPrivateKey(privateKeyFile, []byte{})
	if err != nil {
		FatalError("parse private key file content", err)
	}

	ski := SKI(privateKey)
	fmt.Println(hex.EncodeToString(ski))
}

// SKI returns the subject key identifier of this key.
func SKI(privKey *ecdsa.PrivateKey) []byte {
	if privKey == nil {
		return nil
	}

	// Marshall the public key
	ecdhPk, err := privKey.ECDH()
	if err != nil {
		panic(err)
	}

	raw := ecdhPk.Bytes()
	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func pemToPrivateKey(raw []byte, pwd []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("failed decoding PEM. Block must be different from nil [% x]", raw)
	}

	if x509.IsEncryptedPEMBlock(block) { //nolint:staticcheck
		if len(pwd) == 0 {
			return nil, errors.New("encrypted Key. Need a password")
		}

		decrypted, err := x509.DecryptPEMBlock(block, pwd) //nolint:staticcheck
		if err != nil {
			return nil, fmt.Errorf("failed PEM decryption: %w", err)
		}

		key, err := derToPrivateKey(decrypted)
		if err != nil {
			return nil, err
		}
		return key, err
	}

	key, err := derToPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, err
}

func derToPrivateKey(der []byte) (*ecdsa.PrivateKey, error) {
	if keyi, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch v := keyi.(type) {
		case *ecdsa.PrivateKey:
			return v, nil
		default:
			return nil, errors.New("found unknown private key type in PKCS#8 wrapping")
		}
	}

	return nil, errors.New("invalid key type. The DER must contain an ecdsa.PrivateKey")
}
