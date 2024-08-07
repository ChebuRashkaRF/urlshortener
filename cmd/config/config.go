package config

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	FlagLogLevel    string
}

func NewConfig() *Config {
	parseFlags()

	if serverAddress == "" {
		serverAddress = ":12345"
	}
	if baseURL == "" {
		baseURL = "http://localhost" + serverAddress
	}
	return &Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		FlagLogLevel:    "info",
	}
}

var Cnf *Config
