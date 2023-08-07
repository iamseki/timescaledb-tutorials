package api

import (
	"fmt"
	"time"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/core"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type config struct {
	Name            string `envconfig:"APP_NAME" default:"nyc_taxi_api"`
	LogLevel        string `envconfig:"APP_LOG_LEVEL" default:"info"`
	Port            string `envconfig:"APP_PORT" default:"8080"`
	HTTPReadTimeout int    `envconfig:"HTTP_READ_TIMEOUT" default:"120"`
}

type API struct {
	server *echo.Echo
	logger *zap.Logger
	config *config
	core   *core.Core
}

func New(logger *zap.Logger, core *core.Core) (*API, error) {
	apiConfig := &config{}
	if err := envconfig.Process("cryptocurrency-apy", apiConfig); err != nil {
		logger.Fatal("Error loading API envconfig", zap.Any("apiConfig", apiConfig))
	}

	return &API{
		server: echo.New(),
		logger: logger,
		config: apiConfig,
		core:   core,
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

	api.server.GET("/v1/transactions/volume/last", api.TxVolumeXFeesByHourLast)
	api.server.GET("/v1/transactions/volume/usd/last", api.TxVolumeXUSDByHourLast)

	api.server.GET("/v1/blocks/recent", api.listRecentBlocks)
	api.server.GET("/v1/blocks/mining-fee/transactions/last", api.BlockVolumeTxXMiningFeeByHourLast)
	api.server.GET("/v1/blocks/mining-fee/last", api.BlockVolumeXMiningFeeByHourLast)

	api.server.GET("/v1/mining/reward/last", api.MinersRewardByHourLast)
	api.server.GET("/v1/mining/fees/last", api.MiningBlockWeightXFeeByHourLast)
	api.server.GET("/v1/mining/revenue/last", api.MiningRevenueByHourLast)

	api.logger.Info(fmt.Sprintf("Initializing HTTP server on PORT: %v", api.config.Port))
	return api.server.Start(":" + api.config.Port)
}

func (api *API) Close() error {
	return api.server.Close()
}
