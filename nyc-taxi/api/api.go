package api

import (
	"fmt"
	"time"

	"github.com/iamseki/timescaledb-tutorials/nyc-taxi/repository"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type API struct {
	server     *echo.Echo
	logger     *zap.Logger
	config     *Config
	repository repository.Repository
}

func New(logger *zap.Logger, config *Config, repository repository.Repository) (*API, error) {
	return &API{
		server:     echo.New(),
		logger:     logger,
		config:     config,
		repository: repository,
	}, nil
}

func (api *API) loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		startedAt := time.Now()

		if err := next(ctx); err != nil {
			api.logger.Error("Error on logger middleware handler", zap.Error(err))
			return err
		}

		api.logger.Info("http access logger",
			zap.String("method", ctx.Request().Method),
			zap.Any("header", ctx.Request().Header),
			zap.String("querystring", ctx.QueryString()),
			zap.Int("status_code", ctx.Response().Status),
			zap.Duration("duration", time.Since(startedAt)),
		)
		return nil
	}
}

func (api *API) Start() error {
	api.server.HideBanner = true

	// MIDDLEWARES STUFF
	api.logger.Info("Initializing logger middleware")
	api.server.Use(api.loggerMiddleware)

	// ROUTES
	api.logger.Info("Initializing routes")
	api.server.GET("/v1/healthcheck", api.healthcheck)
	api.server.GET("/v1/rides/day/since", api.ridesByDaySince)

	api.logger.Info(fmt.Sprintf("Initializing HTTP server on PORT: %v", api.config.App.Port))
	return api.server.Start(":" + api.config.App.Port)
}
