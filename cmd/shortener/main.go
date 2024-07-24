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

	logger.Log.Info("Running server on", zap.String("address", cnf.ServerAddress))

	return http.ListenAndServe(cnf.ServerAddress, router.NewRouter())
}

func main() {
	config.Cnf = config.NewConfig()

	db, err := storage.NewDatabase(config.Cnf.DatabaseDSN)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
		panic(err)
	}
	defer db.Close()

	urlStorage, err := storage.NewURLStorage(config.Cnf.FileStoragePath)
	if err != nil {
		logger.Log.Error("Failed to initialize URL storage", zap.Error(err))
		panic(err)
	}

	defer func() {
		if err := urlStorage.Close(); err != nil {
			logger.Log.Error("Error closing URLStorage", zap.Error(err))
		}
	}()

	handler.URLStore = urlStorage
	handler.DB = db

	if err := run(config.Cnf); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}
