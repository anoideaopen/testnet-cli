package cmd

import (
	"crypto/ed25519"
	sha30 "crypto/sha3"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
)

const (
	base58Str      = "base58"
	base58checkStr = "base58check"
	hexStr         = "hex"
	strStr         = "str"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "converter one format to another",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		from := args[0]
		to := args[1]
		convertedString := args[2]
		if from == base58Str && to == hexStr {
			bytes := base58.Decode(convertedString)
			result := hex.EncodeToString(bytes)
			fmt.Println(result)
		} else if from == base58Str && to == "sum3hex" {
			bytes := base58.Decode(convertedString)
			hashed := sha30.Sum256(bytes)
			result := hex.EncodeToString(hashed[:])
			fmt.Println(result)
		} else if from == base58Str && to == "Sum256base58CheckEncode" {
			bytes := base58.Decode(convertedString)
			if len(bytes) != ed25519.PublicKeySize {
				panic(fmt.Sprintf("incorrect publik key size %d\n", len(bytes)))
			}
			hash := sha30.Sum256(bytes)
			result := base58.CheckEncode(hash[1:], hash[0])
			fmt.Println(result)
		} else if from == base58checkStr && to == base58Str {
			bytes, _, err := base58.CheckDecode(convertedString)
			if err != nil {
				FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
			}
			result := base58.Encode(bytes)
			fmt.Println(result)
		} else if from == base58checkStr && to == hexStr {
			bytes, version, err := base58.CheckDecode(convertedString)
			if err != nil {
				FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
			}
			b := make([]byte, 0, len(bytes)+1)
			b = append(b, version)
			b = append(b, bytes...)
			result := hex.EncodeToString(b)
			fmt.Println(result)
		} else if from == base58Str && to == "sum256base58Check" {
			bytes := base58.Decode(convertedString)
			if len(bytes) != ed25519.PublicKeySize {
				fmt.Printf("incorrect public key size %d but expected %d\n", len(bytes), ed25519.PublicKeySize)
			}
			hash := sha30.Sum256(bytes)
			result := base58.CheckEncode(hash[1:], hash[0])
			fmt.Println(result)
		} else if from == base58checkStr && to == "base64" {
			bytes, _, err := base58.CheckDecode(convertedString)
			if err != nil {
				FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
			}
			rawDecodedText, err := base64.StdEncoding.DecodeString(string(bytes))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Decoded text: %s\n", rawDecodedText)
		} else if from == hexStr && to == base58Str {
			bytes, err := hex.DecodeString(convertedString)
			if err != nil {
				FatalError(fmt.Sprintf("Error decode hex from string '%s'", convertedString), err)
			}
			result := base58.Encode(bytes)
			fmt.Println(result)
		} else if from == hexStr && to == base58checkStr {
			bytes, err := hex.DecodeString(convertedString)
			if err != nil {
				FatalError(fmt.Sprintf("Error decode hex from string '%s'", convertedString), err)
			}
			result, err := base58.CheckEncode(bytes[1:], bytes[0]), nil
			if err != nil {
				FatalError(fmt.Sprintf("CheckEncode Error '%s'", convertedString), err)
			}
			fmt.Println(result)
		} else if from == strStr && to == base58Str {
			result := base58.Encode([]byte(convertedString))
			fmt.Println(result)
		} else if from == strStr && to == hexStr {
			result := hex.EncodeToString([]byte(convertedString))
			fmt.Println(result)
		} else if from == strStr && to == base58checkStr {
			result := base58.CheckEncode([]byte(convertedString), 0)
			fmt.Println(result)
		} else if from == base58Str && to == base58checkStr {
			bytes := base58.Decode(convertedString)
			result := base58.CheckEncode(bytes, 0)
			fmt.Println(result)
		} else {
			FatalError(fmt.Sprintf("convert '%s' FROM %s to %s", from, to, convertedString), errors.New("ERROR: unknown type convert"))
		}
	},
}
