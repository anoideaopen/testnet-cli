package config

type ApplicationConfig struct {
	NumberRequest     int    `mapstructure:"numberRequest,omitempty"`
	RequestsPerSecond int    `mapstructure:"requestsPerSecond,omitempty"`
	ChaincodeName     string `mapstructure:"chaincodeName,omitempty"`
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
	PostgresDbName   string `mapstructure:"dbName,omitempty"`
}
