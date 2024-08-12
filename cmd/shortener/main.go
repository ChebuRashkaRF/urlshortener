package main

import (
	"github.com/ChebuRashkaRF/urlshortener/internal/handler"
	"github.com/ChebuRashkaRF/urlshortener/internal/storage"
	"net/http"

	"go.uber.org/zap"

	"github.com/ChebuRashkaRF/urlshortener/cmd/config"
	"github.com/ChebuRashkaRF/urlshortener/internal/logger"
	"github.com/ChebuRashkaRF/urlshortener/internal/router"
)

func run(cnf *config.Config) error {
	if err := logger.Initialize(cnf.FlagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server on", zap.String("address", cnf.ServerAddress), zap.String("storage", string(cnf.StorageType)))

	return http.ListenAndServe(cnf.ServerAddress, router.NewRouter())
}

func main() {
	config.Cnf = config.NewConfig()

	var urlStorage storage.Storage

	switch config.Cnf.StorageType {
	case config.DatabaseStorage:
		db, err := storage.NewDatabaseStorage(config.Cnf.DatabaseDSN)
		if err != nil {
			logger.Log.Fatal("Failed to initialize database storage", zap.Error(err))
			panic(err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				logger.Log.Error("Failed to close database", zap.Error(err))
			}
		}()
		urlStorage = db

	case config.FileStorage:
		fileStorage, err := storage.NewFileStorage(config.Cnf.FileStoragePath)
		if err != nil {
			logger.Log.Fatal("Failed to initialize file storage", zap.Error(err))
			panic(err)
		}
		defer func() {
			if err := fileStorage.Close(); err != nil {
				logger.Log.Error("Error closing file storage", zap.Error(err))
			}
		}()
		urlStorage = fileStorage
	case config.MemoryStorage:
		urlStorage = storage.NewMemoryStorage()
	default:
		logger.Log.Fatal("Invalid storage type")
		panic("Invalid storage type")
	}

	handler.URLStore = urlStorage

	if err := run(config.Cnf); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
