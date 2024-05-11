package config

type Config struct {
	ServerAddress string
	BaseURL       string
}

func NewConfig(serverAddress, baseURL string) *Config {
	if serverAddress == "" {
		serverAddress = ":12345"
	}
	if baseURL == "" {
		baseURL = "http://localhost" + serverAddress
	}
	return &Config{
		ServerAddress: serverAddress,
		BaseURL:       baseURL,
	}
}

var Cnf *Config
