package config

import (
	"fmt"

	"github.com/anoideaopen/foundation/proto"
)

type ApplicationConfig struct {
	NumberRequest     int         `mapstructure:"numberRequest,omitempty"`
	RequestsPerSecond int         `mapstructure:"requestsPerSecond,omitempty"`
	ChaincodeName     string      `mapstructure:"chaincodeName,omitempty"`
	KeyType           KeyTypeFlag `mapstructure:"keyType,omitempty"` // ed25519, secp256k1, gost
	// key for sign transaction arguments in hyperledger
	SecretKey          string `mapstructure:"secretKey,omitempty"`
	PrivateKeyFilePath string `mapstructure:"privateKeyFilePath,omitempty"`
	Peers              string `mapstructure:"peers,omitempty"`
	ResponseType       string `mapstructure:"responseType,omitempty"`
	WaitBatch          bool   `mapstructure:"waitBatch,omitempty"`
	// connection with hlf configuration
	Connection   string `mapstructure:"connection,omitempty"`
	Organization string `mapstructure:"organization,omitempty"`
	User         string `mapstructure:"user,omitempty"`
	Observer     ObserverConfig
	Postgres     PostgresConfig
}

// ObserverConfig connection with Observer API
type ObserverConfig struct {
	ObserverUsername string `mapstructure:"username,omitempty"`
	ObserverPassword string `mapstructure:"password,omitempty"`
	ObserverURL      string `mapstructure:"url,omitempty"`
	ObserverVersion  string `mapstructure:"version,omitempty"`
}

type PostgresConfig struct {
	PostgresHost     string `mapstructure:"host,omitempty"`
	PostgresPort     string `mapstructure:"port,omitempty"`
	PostgresUser     string `mapstructure:"user,omitempty"`
	PostgresPassword string `mapstructure:"password,omitempty"`
	PostgresDBName   string `mapstructure:"dbName,omitempty"`
}

type KeyTypeFlag proto.KeyType

func (k *KeyTypeFlag) Set(s string) error {
	fmt.Printf("Set called with s=%q\n", s)
	v, ok := proto.KeyType_value[s]
	if !ok {
		return fmt.Errorf("invalid key type: %s, must be one of %v", s, proto.KeyType_value)
	}
	*k = KeyTypeFlag(v)
	return nil
}

func (k *KeyTypeFlag) String() string {
	if name, ok := proto.KeyType_name[int32(*k)]; ok {
		return name
	}
	return ""
}

func (k *KeyTypeFlag) Type() string {
	return "KeyType"
}

func (k *KeyTypeFlag) Value() proto.KeyType {
	return proto.KeyType(*k)
}
