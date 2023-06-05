package main

import (
	"log"

	"github.com/iamseki/timescaledb-tutorials/nyc-taxi/api"
	"github.com/iamseki/timescaledb-tutorials/nyc-taxi/repository"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	var apiConfig api.Config
	if err := envconfig.Process("nyc-taxi-app", &apiConfig); err != nil {
		logger.Fatal("Error loading API envconfig", zap.Any("apiConfig", apiConfig))
	}

	var repoConfig repository.Config
	if err := envconfig.Process("nyc-taxi-app", &repoConfig); err != nil {
		logger.Fatal("Error loading REPOSITORY envconfig", zap.Any("repoConfig", repoConfig))
	}

	repository, err := repository.NewPostgreSQL(&repoConfig, logger)
	if err != nil {
		logger.Fatal("error on repository.NewPostgreSQL", zap.Error(err))
	}

	app, err := api.New(logger, &apiConfig, repository)
	if err != nil {
		logger.Fatal("error on api.New", zap.Error(err))
	}

	app.Start()
}
