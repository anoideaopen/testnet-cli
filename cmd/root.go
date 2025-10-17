package cmd

import (
	"fmt"
	"os"
	"strings"

	config2 "github.com/anoideaopen/testnet-cli/config"
	"github.com/anoideaopen/testnet-cli/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	ENVPREFX = "CLI"
)

var configFilePath string

var (
	config  config2.ApplicationConfig
	rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "cli request to hlf blockchain for atmz network",
		Long:  `cli request to hlf blockchain for atmz network`,
	}
)

// Execute executes the root command.
func Execute(version string, commit string, date string) error {
	Version = version
	Commit = commit
	Date = date

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "the config file path")
	rootCmd.PersistentFlags().IntVarP(&config.NumberRequest, "numberRequest", "n", 1, "number request")
	rootCmd.PersistentFlags().IntVarP(&config.RequestsPerSecond, "requestsPerSecond", "t", 1, "requests per second")
	rootCmd.PersistentFlags().StringVarP(&config.ChaincodeName, "chaincode", "c", "", "chaincode name on which this command should be executed")
	rootCmd.PersistentFlags().VarP(&config.KeyType, "keyType", "k", "key type: ed25519 - default | secp256k1 | gost")
	rootCmd.PersistentFlags().StringVarP(&config.SecretKey, "secretKey", "s", "", "private key in format base58 or base58check or hex (you can use func 'cli privkey' to gen private key)")
	rootCmd.PersistentFlags().StringVarP(&config.Peers, "peers", "p", "", "parameter for invoke request. don't wait event 'batchExecute'")
	rootCmd.PersistentFlags().StringVarP(&config.ResponseType, "responseType", "r", "resp", "response type 'tx','resp'")
	rootCmd.PersistentFlags().BoolVarP(&config.WaitBatch, "waitBatch", "w", false, "wait batch executed by robot")

	// connection with hlf configuration
	rootCmd.PersistentFlags().StringVarP(&config.Connection, "connection", "f", "", "path to connection.yaml. The path to fabric sdk config file for connection to HLF")
	rootCmd.PersistentFlags().StringVarP(&config.Organization, "organization", "", "testnet", "organization in connection.yaml. Organization name where store user for FabricSDK 'Username'")
	rootCmd.PersistentFlags().StringVarP(&config.User, "username", "u", "backend", "user name. User name for sign request for FabricSDK")

	rootCmd.PersistentFlags().StringVar(&config.PrivateKeyFilePath, "privateKeyFilePath", "", "private key in file \n-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----")

	rootCmd.PersistentFlags().StringVar(&config.Observer.ObserverUsername, "observerUsername", "", "observerUsername")
	rootCmd.PersistentFlags().StringVar(&config.Observer.ObserverPassword, "observerPassword", "", "observerPassword")
	rootCmd.PersistentFlags().StringVar(&config.Observer.ObserverURL, "observerURL", "", "observerURL")
	rootCmd.PersistentFlags().StringVar(&config.Observer.ObserverVersion, "observerVersion", "v2", "observerVersion")
	rootCmd.PersistentFlags().StringVar(&config.Postgres.PostgresHost, "postgresHost", "", "postgresHost")
	rootCmd.PersistentFlags().StringVar(&config.Postgres.PostgresPort, "postgresPort", "", "postgresPort")
	rootCmd.PersistentFlags().StringVar(&config.Postgres.PostgresUser, "postgresUser", "", "postgresUser")
	rootCmd.PersistentFlags().StringVar(&config.Postgres.PostgresPassword, "postgresPassword", "", "postgresPassword")
	rootCmd.PersistentFlags().StringVar(&config.Postgres.PostgresDBName, "postgresDbName", "", "postgresDbName")

	rootCmd.AddCommand(privkeyCmd)
	rootCmd.AddCommand(pubkeyCmd)
	rootCmd.AddCommand(addressCmd)
	rootCmd.AddCommand(skiCmd)

	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(invokeCmd)
	rootCmd.AddCommand(scriptCmd)

	rootCmd.AddCommand(blockByIDCmd)
	rootCmd.AddCommand(channelHeightCmd)
	rootCmd.AddCommand(txCmd)

	rootCmd.AddCommand(statusCmd)
	// rootCmd.AddCommand(fetchBatchCmd)

	// rootCmd.AddCommand(invokeACLCmd)
	// rootCmd.AddCommand(chaincodeVersionCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(getTxIDFromBlockCmd)
	rootCmd.AddCommand(batchExecuteCmd)
	// rootCmd.AddCommand(validateBlockCmd)
	// rootCmd.AddCommand(generateMessageCmd)
	// rootCmd.AddCommand(sendRequestCmd)
	// rootCmd.AddCommand(signMessageCmd)
	// // service discovery request
	// rootCmd.AddCommand(discoveryQueryPeerQueryCmd)
	// rootCmd.AddCommand(discoveryQueryCcQueryCmd)
}

func initConfig() {
	isHLFMethods := os.Args[1] == "invoke" ||
		os.Args[1] == "query" ||
		os.Args[1] == "block" ||
		os.Args[1] == "channelHeight"
	if isHLFMethods && configFilePath == "" {
		panic("config file can't be empty")
	}
	if isHLFMethods && pathNotExist(configFilePath) {
		panic(fmt.Sprintf("file path %s doesn't exists", configFilePath))
	}

	// Set the config file path.
	viper.SetConfigFile(configFilePath)

	// Read the config file.
	err := viper.ReadInConfig()
	if err != nil {
		logger.GetLogger().Warn("config file not found", zap.Error(err))
	}

	viper.SetEnvPrefix(ENVPREFX)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		FatalError("failed unmarshal config", err)
	}
}

func pathNotExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}
