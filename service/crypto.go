package service

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/anoideaopen/foundation/keys"
	"github.com/anoideaopen/foundation/proto"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

// Sign creates a signed message using the provided key and transaction details.
// It constructs a message from the method name, chaincode, channel, arguments, nonce, and public key.
// The message is then signed using the key's type (Ed25519, secp256k1, etc.).
// Returns the message with signature, the message hash, and an error if any occurred.
func Sign(k *keys.Keys, channel string, chaincode string, methodName string, args []string) ([]string, string, error) {
	nonce := GetNonce()
	publicKeyBase58, err := ConvertPublicKeyToBase58(k)
	if err != nil {
		return nil, "", err
	}
	result := append(append([]string{methodName, uuid.NewString(), chaincode, channel}, args...), nonce, publicKeyBase58)
	m := []byte(strings.Join(result, ""))

	logger.Debug(
		"For sign",
		zap.Strings("result", result),
	)

	message, signature, err := keys.SignMessageByKeyType(k.KeyType, k, m)
	if err != nil {
		return nil, "", err
	}

	var messageWithSig []string
	messageWithSig = append(append(messageWithSig, result[1:]...), base58.Encode(signature))
	hash := hex.EncodeToString(message)

	logger.Debug(
		"Sign result",
		zap.Strings("messageWithSig", messageWithSig),
		zap.String("hash", hash),
	)

	return messageWithSig, hash, nil
}

// SignACL creates a multi-signature message for ACL operations.
// It builds a message from the provided method name, address, reason, reasonID, new public key, and nonce.
// Each signer signs the same message using their private key, and all signatures are collected.
// Returns the message with all signatures, the message hash, and an error if any occurred.
func SignACL(signers []*keys.Keys, methodName string, address string, reason string, reasonID string, newPkey string) ([]string, string, error) {
	nonce := GetNonce()

	result := []string{methodName, address, reason, reasonID, newPkey, nonce}
	for _, k := range signers {
		pubBase58, err := ConvertPublicKeyToBase58(k)
		if err != nil {
			return nil, "", err
		}
		result = append(result, pubBase58)
	}

	logger.Debug("For sign", zap.Strings("result", result))
	message := []byte(strings.Join(result, ""))

	signatures := make([]string, 0)
	for _, k := range signers {
		_, signature, err := keys.SignMessageByKeyType(k.KeyType, k, message)
		if err != nil {
			return nil, "", err
		}

		signatures = append(signatures, hex.EncodeToString(signature))
	}

	var messageWithSig []string
	messageWithSig = append(append(messageWithSig, result[1:]...), signatures...)
	hash := hex.EncodeToString(message)

	logger.Debug(
		"Sign result",
		zap.Strings("messageWithSig", messageWithSig),
		zap.String("hash", hash),
	)

	return messageWithSig, hash, nil
}

// SignMessage creates a message from the provided `result` slice and signs it using the specified key type.
// It supports multiple key algorithms (e.g., Ed25519, secp256k1) through `keys.SignMessageByKeyType`.
// Returns the signature, the original message bytes, and an error if signing fails.
func SignMessage(k *keys.Keys, keyType proto.KeyType, result []string) ([]byte, []byte, error) {
	m := []byte(strings.Join(result, ""))

	// Подписываем сообщение в зависимости от типа ключа
	message, signature, err := keys.SignMessageByKeyType(keyType, k, m)
	if err != nil {
		logger.Error("SignMessageByKeyType", zap.Error(err))
		return nil, message, err
	}

	return signature, message, nil
}

// --- Get Addresses ---

// GetAddress - get address by encoded string in standard encoded for project is 'base58.Check'
// secretKey string - private key in base58check, or hex or base58
func GetAddress(secretKey string, keyType proto.KeyType) (string, error) {
	publicKey, err := GetPublicKey(secretKey, keyType)
	if err != nil {
		return "", err
	}

	return GetAddressByPublicKey(publicKey)
}

// GetAddressByPublicKey - get address by encoded string in standard encoded for project is 'base58.Check'
// secretKey string - private key in base58check, or hex or base58
func GetAddressByPublicKey(publicKey string) (string, error) {
	if len(publicKey) == 0 {
		return "", errors.New("publicKey can't be empty")
	}

	pubBytes := base58.Decode(publicKey)
	if len(pubBytes) == 0 {
		return "", errors.New("decoded public key is empty")
	}

	hash := sha3.Sum256(pubBytes)
	return base58.CheckEncode(hash[1:], hash[0]), nil
}

// --- Generate ---

// GeneratePrivateKey generates a new private key for the given key type
// and returns it encoded in base58check format.
func GeneratePrivateKey(keyType proto.KeyType) (string, error) {
	k, err := GeneratePrivateAndPublicKey(keyType)
	if err != nil {
		return "", err
	}

	var b []byte
	switch keyType {
	case proto.KeyType_ed25519:
		b = k.PrivateKeyEd25519
	case proto.KeyType_secp256k1:
		b = crypto.FromECDSA(k.PrivateKeySecp256k1)
	case proto.KeyType_gost:
		return "", errors.New("GOST key generation not implemented")
	default:
		return "", errors.New("unsupported key type")
	}

	return ConvertPrivateKeyToBase58CheckFromBytes(b), nil
}

// GeneratePrivateAndPublicKey generates a key pair based on the provided key type.
// Supports Ed25519 and secp256k1 (Ethereum). For GOST keys, returns a not implemented error.
func GeneratePrivateAndPublicKey(keyType proto.KeyType) (*keys.Keys, error) {
	k := &keys.Keys{
		KeyType: keyType,
	}

	switch keyType {
	case proto.KeyType_ed25519:
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		k.KeyType = proto.KeyType_ed25519
		k.PublicKeyEd25519 = publicKey
		k.PrivateKeyEd25519 = privateKey
		return k, nil

	case proto.KeyType_secp256k1:
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			return nil, err
		}
		k.KeyType = proto.KeyType_secp256k1
		publicKey := &privateKey.PublicKey
		k.PublicKeySecp256k1 = publicKey
		k.PrivateKeySecp256k1 = privateKey
		return k, nil

	case proto.KeyType_gost:
		return nil, errors.New("GOST key generation not implemented yet")

	default:
		return nil, errors.New("unsupported key type")
	}
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

// --- Convert ---

// ConvertPrivateKeyToHex converts a private key to a hex string based on the key type.
// Supports Ed25519 and secp256k1. Returns an error for unsupported key types.
func ConvertPrivateKeyToHex(k *keys.Keys) (string, error) {
	switch k.KeyType {
	case proto.KeyType_ed25519:
		if len(k.PrivateKeyEd25519) == 0 {
			return "", errors.New("ed25519 private key is empty")
		}
		return hex.EncodeToString(k.PrivateKeyEd25519), nil

	case proto.KeyType_secp256k1:
		if k.PrivateKeySecp256k1 == nil {
			return "", errors.New("secp256k1 private key is nil")
		}
		privBytes := crypto.FromECDSA(k.PrivateKeySecp256k1)
		return hex.EncodeToString(privBytes), nil

	case proto.KeyType_gost:
		return "", errors.New("GOST key type not implemented")

	default:
		return "", errors.New("unsupported key type: " + k.KeyType.String())
	}
}

// ConvertPrivateKeyToBase58Check - use privateKey with standard encoded type - Base58Check
func ConvertPrivateKeyToBase58CheckFromBytes(privateKey []byte) string {
	encoded := base58.CheckEncode(privateKey[1:], privateKey[0])
	return encoded
}

// ConvertPrivateKeyToBase58Check converts the private key to a Base58Check string
// based on the key type. Supports Ed25519 and secp256k1. Returns an error for unsupported key types.
func ConvertPrivateKeyToBase58Check(k *keys.Keys) (string, error) {
	var privateKeyBytes []byte

	switch k.KeyType {
	case proto.KeyType_ed25519:
		if len(k.PrivateKeyEd25519) == 0 {
			return "", errors.New("ed25519 private key is empty")
		}
		privateKeyBytes = k.PrivateKeyEd25519

	case proto.KeyType_secp256k1:
		if k.PrivateKeySecp256k1 == nil {
			return "", errors.New("secp256k1 private key is nil")
		}
		privateKeyBytes = crypto.FromECDSA(k.PrivateKeySecp256k1)

	case proto.KeyType_gost:
		return "", errors.New("GOST key type not implemented")

	default:
		return "", errors.New("unsupported key type: " + k.KeyType.String())
	}

	encoded := base58.CheckEncode(privateKeyBytes[1:], privateKeyBytes[0])
	return encoded, nil
}

// ConvertPublicKeyToBase58 - use publicKey with standard encoded type - Base58
func ConvertPublicKeyToBase58(k *keys.Keys) (string, error) {
	switch k.KeyType {
	case proto.KeyType_ed25519:
		if k.PublicKeyEd25519 == nil {
			return "", errors.New("ed25519 public key is nil")
		}
		k.PublicKeyBase58 = base58.Encode(k.PublicKeyEd25519)
		return k.PublicKeyBase58, nil

	case proto.KeyType_secp256k1:
		if k.PublicKeySecp256k1 == nil {
			return "", errors.New("secp256k1 public key is nil")
		}
		pubBytes, err := x509.MarshalPKIXPublicKey(k.PublicKeySecp256k1)
		if err != nil {
			return "", errors.New("failed to marshal secp256k1 public key: " + err.Error())
		}
		k.PublicKeyBase58 = base58.Encode(pubBytes)
		return k.PublicKeyBase58, nil

	case proto.KeyType_gost:
		// TODO: реализовать для GOST позже
		return "", errors.New("GOST key type not yet implemented")
	default:
		return "", errors.New("unsupported key type: " + k.KeyType.String())
	}
}

// ConvertSignatureToBase58 - use signature with standard encoded type - Base58
func ConvertSignatureToBase58(publicKey []byte) string {
	encoded := base58.Encode(publicKey)
	return encoded
}

// --- Get Keys ---

// GetPublicKey in standard encoded for project is 'base58'
// secretKey string - private key in base58check, or hex or base58
func GetPublicKey(secretKey string, keyType proto.KeyType) (string, error) {
	k, err := GetKeys(secretKey, keyType)
	if err != nil {
		return "", err
	}

	switch k.KeyType {
	case proto.KeyType_ed25519:
		if k.PublicKeyEd25519 == nil {
			return "", errors.New("ed25519 public key is nil")
		}
		return base58.Encode(k.PublicKeyEd25519), nil

	case proto.KeyType_secp256k1:
		if k.PublicKeySecp256k1 == nil {
			return "", errors.New("secp256k1 public key is nil")
		}
		pubBytes, err := k.PublicKeySecp256k1.Bytes()
		if err != nil {
			return "", err
		}
		return base58.Encode(pubBytes), nil

	case proto.KeyType_gost:
		return "", errors.New("GOST key type decoding is not yet implemented")

	default:
		return "", fmt.Errorf("unsupported key type: %v", k.KeyType)
	}
}

// GetKeys - get private key type Ed25519 by encoded private key in string
// secretKey string - private key in base58check, or hex or base58
func GetKeys(secretKey string, keyType proto.KeyType) (*keys.Keys, error) {
	switch keyType {
	case proto.KeyType_ed25519:
		privateKey, publicKey, err := GetPrivateKeySKFromBase58Check(secretKey)
		if err != nil {
			privateKey, publicKey, err = GetPrivateKeySKFromHex(secretKey)
			if err != nil {
				privateKey, publicKey, err = GetPrivateKeySKFromBase58(secretKey)
				if err != nil {
					return nil, err
				}
			}
		}
		return &keys.Keys{
			KeyType:           proto.KeyType_ed25519,
			PublicKeyEd25519:  publicKey,
			PrivateKeyEd25519: privateKey,
		}, nil
	case proto.KeyType_secp256k1:
		privateKey, publicKey, err := GetSecp256k1KeysFromBase58Check(secretKey)
		if err != nil {
			privateKey, publicKey, err = GetSecp256k1KeysFromHex(secretKey)
			if err != nil {
				privateKey, publicKey, err = GetSecp256k1KeysFromBase58(secretKey)
				if err != nil {
					return nil, err
				}
			}
		}
		return &keys.Keys{
			KeyType:             proto.KeyType_secp256k1,
			PublicKeySecp256k1:  publicKey,
			PrivateKeySecp256k1: privateKey,
		}, nil
	case proto.KeyType_gost:
		return nil, errors.New("GOST key type decoding is not yet implemented")

	default:
		return nil, errors.New("unsupported key type: " + keyType.String())
	}
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

// GetSecp256k1KeysFromBase58Check - get secp256k1 private key by Base58Check encoded string
func GetSecp256k1KeysFromBase58Check(secretKey string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	decoded, _, err := base58.CheckDecode(secretKey)
	if err != nil {
		return nil, nil, err
	}
	priv, err := crypto.ToECDSA(decoded)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// GetSecp256k1KeysFromHex - get secp256k1 private key by hex encoded string
func GetSecp256k1KeysFromHex(secretKey string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privBytes, err := hex.DecodeString(secretKey)
	if err != nil {
		return nil, nil, err
	}
	priv, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// GetSecp256k1KeysFromBase58 - get secp256k1 private key by Base58 encoded string
func GetSecp256k1KeysFromBase58(secretKey string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privBytes := base58.Decode(secretKey)
	priv, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}
