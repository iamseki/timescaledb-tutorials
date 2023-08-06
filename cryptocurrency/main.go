package main

import (
	"context"
	"log"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/api"
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/core"
	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/repository"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewHTTPServer(lc fx.Lifecycle, api *api.API, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := api.Start()
				if err != nil {
					logger.Sugar().Fatalln("Error on start api:", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return api.Close()
		},
	})
}

func FxLogger(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}

func NewAppLogger() *zap.Logger {
	zapCfg := zap.NewProductionConfig()
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

func main() {
	fx.New(
		fx.Provide( // Provide must receives every constructors and related dependencies, the fx will resolve them all
			NewAppLogger,
			api.New,
			core.New,
			fx.Annotate(
				repository.NewPostgreSQL,
				fx.As(new(repository.Repository)),
			),
		),
		fx.WithLogger(FxLogger),  // overriding the default logger
		fx.Invoke(NewHTTPServer), // using the fx life cycle to controls gracefull shutdown and app startup as well
	).Run()
}
