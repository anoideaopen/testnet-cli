package cmd

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/sha3"
)

var (
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "converter one format to another",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			from := args[0]
			to := args[1]
			convertedString := args[2]

			if from == "base58" && to == "hex" {
				bytes := base58.Decode(convertedString)
				result := hex.EncodeToString(bytes)
				fmt.Println(result)
			} else if from == "base58" && to == "sum3hex" {
				bytes := base58.Decode(convertedString)
				hashed := sha3.Sum256(bytes)
				result := hex.EncodeToString(hashed[:])
				fmt.Println(result)
			} else if from == "base58" && to == "Sum256base58CheckEncode" {
				bytes := base58.Decode(convertedString)
				if len(bytes) != ed25519.PublicKeySize {
					panic(fmt.Sprintf("incorrect publik key size %d\n", len(bytes)))
				}
				hash := sha3.Sum256(bytes)
				result := base58.CheckEncode(hash[1:], hash[0])
				fmt.Println(result)
			} else if from == "base58check" && to == "base58" {
				bytes, _, err := base58.CheckDecode(convertedString)
				if err != nil {
					FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
				}
				result := base58.Encode(bytes)
				fmt.Println(result)
			} else if from == "base58check" && to == "hex" {
				bytes, version, err := base58.CheckDecode(convertedString)
				if err != nil {
					FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
				}
				b := []byte{version}
				b = append(b, bytes...)
				result := hex.EncodeToString(b)
				fmt.Println(result)
			} else if from == "base58" && to == "sum256base58Check" {
				bytes := base58.Decode(convertedString)
				if len(bytes) != ed25519.PublicKeySize {
					fmt.Printf("incorrect publik key size %d but expected %d\n", len(bytes), ed25519.PublicKeySize)
				}
				hash := sha3.Sum256(bytes)
				result := base58.CheckEncode(hash[1:], hash[0])
				fmt.Println(result)
			} else if from == "base58check" && to == "base64" {
				bytes, _, err := base58.CheckDecode(convertedString)
				if err != nil {
					FatalError(fmt.Sprintf("Error decode base58check from string '%s'", convertedString), err)
				}
				rawDecodedText, err := base64.StdEncoding.DecodeString(string(bytes))
				if err != nil {
					panic(err)
				}
				fmt.Printf("Decoded text: %s\n", rawDecodedText)
			} else if from == "hex" && to == "base58" {
				bytes, err := hex.DecodeString(convertedString)
				if err != nil {
					FatalError(fmt.Sprintf("Error decode hex from string '%s'", convertedString), err)
				}
				result := base58.Encode(bytes)
				fmt.Println(result)
			} else if from == "hex" && to == "base58check" {
				bytes, err := hex.DecodeString(convertedString)
				if err != nil {
					FatalError(fmt.Sprintf("Error decode hex from string '%s'", convertedString), err)
				}
				result, err := base58.CheckEncode(bytes[1:], bytes[0]), nil
				if err != nil {
					FatalError(fmt.Sprintf("CheckEncode Error '%s'", convertedString), err)
				}
				fmt.Println(result)
			} else if from == "str" && to == "base58" {
				result := base58.Encode([]byte(convertedString))
				fmt.Println(result)
			} else if from == "str" && to == "hex" {
				result := hex.EncodeToString([]byte(convertedString))
				fmt.Println(result)
			} else if from == "str" && to == "base58check" {
				result := base58.CheckEncode([]byte(convertedString), 0)
				fmt.Println(result)
			} else if from == "base58" && to == "base58check" {
				bytes := base58.Decode(convertedString)
				result := base58.CheckEncode(bytes, 0)
				fmt.Println(result)
			} else {
				FatalError(fmt.Sprintf("convert '%s' FROM %s to %s", from, to, convertedString), errors.New("ERROR: unknown type convert"))
			}
		},
	}
)
