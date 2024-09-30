package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/btcsuite/btcutil/base58"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
)

type SignerInfo struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

// GetPublicKey in standard encoded for project is 'base58'
// secretKey string - private key in base58check, or hex or base58
func GetPublicKey(secretKey string) (string, error) {
	var publicKey ed25519.PublicKey
	var err error

	_, publicKey, err = GetPrivateKey(secretKey)
	if err != nil {
		return "", err
	}

	return base58.Encode(publicKey), nil
}

// GetPrivateKey - get private key type Ed25519 by encoded private key in string
// secretKey string - private key in base58check, or hex or base58
func GetPrivateKey(secretKey string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privateKey, publicKey, err := GetPrivateKeySKFromBase58Check(secretKey)
	if err != nil {
		privateKey, publicKey, err = GetPrivateKeySKFromHex(secretKey)
		if err != nil {
			privateKey, publicKey, err = GetPrivateKeySKFromBase58(secretKey)
		}
	}

	return privateKey, publicKey, err
}

// GetPrivateKeySKFromBase58Check - get private key type Ed25519 by string - Base58Check encoded private key
// secretKey string - private key in Base58Check
func GetPrivateKeySKFromBase58Check(secretKey string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	decode, ver, err := base58.CheckDecode(secretKey)
	if err != nil {
		return nil, nil, err
	}
	privateKey := ed25519.PrivateKey(append([]byte{ver}, decode...))
	pubKey, _ := privateKey.Public().(ed25519.PublicKey)
	return privateKey, pubKey, nil
}

// GetPrivateKeySKFromHex - get private key type Ed25519 by string - hex encoded private key
// secretKey string - private key in hex
func GetPrivateKeySKFromHex(secretKey string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privateKey, err := hex.DecodeString(secretKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, _ := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
	return privateKey, pubKey, nil
}

// GetPrivateKeySKFromBase58 - get private key type Ed25519 by string - Base58 encoded private key
// secretKey string - private key in Base58
func GetPrivateKeySKFromBase58(secretKey string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	privateKey := base58.Decode(secretKey)
	pubKey, _ := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
	return privateKey, pubKey, nil
}

func GeneratePrivateKey() (string, error) {
	_, privateKey, err := GeneratePrivateAndPublicKey()
	if err != nil {
		return "", err
	}

	return ConvertPrivateKeyToBase58Check(privateKey), nil
}

func GeneratePrivateAndPublicKey() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	return publicKey, privateKey, err
}

// ConvertPrivateKeyToBase58Check - use privateKey with standard encoded type - Base58Check
func ConvertPrivateKeyToBase58Check(privateKey ed25519.PrivateKey) string {
	hash := []byte(privateKey)
	encoded := base58.CheckEncode(hash[1:], hash[0])
	return encoded
}

// ConvertPrivateKeyToHex - use privateKey with standard encoded type - hex
func ConvertPrivateKeyToHex(privateKey ed25519.PrivateKey) string {
	return hex.EncodeToString([]byte(privateKey))
}

func SignACL(signerInfoArray []SignerInfo, methodName string, address string, reason string, reasonID string, newPkey string) ([]string, string, error) {
	nonce := GetNonce()
	// 1. update to change any transactions
	// 2.
	result := []string{methodName, address, reason, reasonID, newPkey, nonce}
	for _, signerInfo := range signerInfoArray {
		result = append(result, ConvertPublicKeyToBase58(signerInfo.PublicKey))
	}

	logger.Debug(
		"For sign",
		zap.Strings("result", result),
	)

	message := sha3.Sum256([]byte(strings.Join(result, "")))

	signatures := make([]string, 0)
	for _, signerInfo := range signerInfoArray {
		sig := ed25519.Sign(signerInfo.PrivateKey, message[:])
		if !ed25519.Verify(signerInfo.PublicKey, message[:], sig) {
			err := errors.New("valid signature rejected")
			logger.Error("ed25519.Verify", zap.Error(err))
			return nil, "", err
		}
		signatures = append(signatures, hex.EncodeToString(sig))
	}

	messageWithSig := result[1:]
	messageWithSig = append(messageWithSig, signatures...)
	hash := hex.EncodeToString(message[:])

	logger.Debug(
		"Sign result",
		zap.Strings("messageWithSig", messageWithSig),
		zap.String("hash", hash),
	)

	return messageWithSig, hash, nil
}

func GetNonce() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// ConvertPublicKeyToBase58 - use publicKey with standard encoded type - Base58
func ConvertPublicKeyToBase58(publicKey ed25519.PublicKey) string {
	encoded := base58.Encode(publicKey)
	return encoded
}

// GetAddress - get address by encoded string in standard encoded for project is 'base58.Check'
// secretKey string - private key in base58check, or hex or base58
func GetAddress(secretKey string) (string, error) {
	var publicKey ed25519.PublicKey
	var err error

	_, publicKey, err = GetPrivateKey(secretKey)
	if err != nil {
		return "", err
	}

	return GetAddressByPublicKey(publicKey)
}

// GetAddressByPublicKey - get address by encoded string in standard encoded for project is 'base58.Check'
// secretKey string - private key in base58check, or hex or base58
func GetAddressByPublicKey(publicKey ed25519.PublicKey) (string, error) {
	if len(publicKey) == 0 {
		return "", errors.New("publicKey can't be empty")
	}

	hash := sha3.Sum256(publicKey)
	return base58.CheckEncode(hash[1:], hash[0]), nil
}

func SignMessage(signerInfo SignerInfo, result []string) ([]byte, [32]byte, error) {
	message := sha3.Sum256([]byte(strings.Join(result, "")))
	sig := ed25519.Sign(signerInfo.PrivateKey, message[:])
	if !ed25519.Verify(signerInfo.PublicKey, message[:], sig) {
		err := errors.New("valid signature rejected")
		logger.Error("ed25519.Verify", zap.Error(err))
		return nil, message, err
	}
	return sig, message, nil
}

func GenerateMessage(validatorPublicKeys []string, channelID string, chaincodeName string, methodName string, args []string) string {
	requestID := ""
	nonce := GetNonce()
	result := append(append([]string{methodName, requestID, chaincodeName, channelID}, args...), nonce)
	result = append(result, validatorPublicKeys...)

	logger.Debug(
		"For sign",
		zap.Strings("result", result),
	)

	return strings.Join(result, "\n")
}
