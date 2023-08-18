package api

import (
	"net/http"

	"github.com/iamseki/timescaledb-tutorials/cryptocurrency/core"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) insertTransaction(c echo.Context) error {
	var body core.TransactionInput
	if err := c.Bind(&body); err != nil {
		api.logger.Error("Error on parse json body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	api.logger.Info("Inserting transaction", zap.Any("transactionInput", body))
	err := api.core.InsertTransaction(body)
	if err != nil {
		api.logger.Error("Error calling core.InsertTransaction", zap.Any("body", body))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
