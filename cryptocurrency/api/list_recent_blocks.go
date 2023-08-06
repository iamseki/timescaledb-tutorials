package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) listRecentBlocks(c echo.Context) error {
	limitStr := c.QueryParam("limit")

	if limitStr == "" {
		api.logger.Error("limit cannot be empty", zap.String("limit", limitStr), zap.String("requestURI", c.Request().RequestURI))
		return echo.NewHTTPError(http.StatusBadRequest, "limit cannot be empty")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		api.logger.Error("limit must be a number", zap.Error(err), zap.Any("limit", limit))
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be a number")
	}

	res, err := api.core.ListRecentBlocks(limit)
	if err != nil {
		api.logger.Error("Error on ListRecentBlocks", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &res)
}
