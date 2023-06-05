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

	var config api.Config

	if err := envconfig.Process("nyc-taxi-app", &config); err != nil {
		logger.Fatal("Error loading envconfig", zap.Any("config", config))
	}

	db, err := repository.NewPostgreSQL(&config, logger)
	if err != nil {
		logger.Fatal("error on repository.NewPostgreSQL", zap.Error(err))
	}

	db.Close()

	app, err := api.New(logger, &config)
	if err != nil {
		logger.Fatal("error on api.New", zap.Error(err))
	}

	app.Start()
}
