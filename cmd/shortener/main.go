package main

import (
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

	if err := run(config.Cnf); err != nil {
		panic(err)
	}
}
