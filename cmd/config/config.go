package config

type StorageType string

const (
	DatabaseStorage StorageType = "db"
	FileStorage     StorageType = "file"
	MemoryStorage   StorageType = "memory"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	FlagLogLevel    string
	DatabaseDSN     string
	StorageType     StorageType
}

func NewConfig() *Config {
	parseFlags()

	if serverAddress == "" {
		serverAddress = ":12345"
	}
	if baseURL == "" {
		baseURL = "http://localhost" + serverAddress
	}

	storageType := determineStorageType()

	return &Config{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL,
		FileStoragePath: fileStoragePath,
		FlagLogLevel:    "info",
		DatabaseDSN:     databaseDSN,
		StorageType:     storageType,
	}
}

var Cnf *Config

func determineStorageType() StorageType {
	if databaseDSN != "" {
		return DatabaseStorage
	} else if fileStoragePath != "" {
		return FileStorage
	} else {
		return MemoryStorage
	}
}
