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
		baseURL = "http://localhost:12345"
	}
	return &Config{
		ServerAddress: serverAddress,
		BaseURL:       baseURL,
	}
}

var Cnf *Config
