package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) MinersRewardByHourLast(c echo.Context) error {
	var days int
	var err error
	daysStr := c.QueryParam("days")

	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			api.logger.Error("days must be a number", zap.Error(err), zap.Any("days", daysStr))
			return echo.NewHTTPError(http.StatusBadRequest, "days must be a number")
		}
	}

	res, err := api.core.MinersRewardByHour(days)
	if err != nil {
		api.logger.Error("Error on core.MinersRewardByHour", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &res)
}

func (api *API) MiningBlockWeightXFeeByHourLast(c echo.Context) error {
	var days int
	var err error
	daysStr := c.QueryParam("days")

	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			api.logger.Error("days must be a number", zap.Error(err), zap.Any("days", daysStr))
			return echo.NewHTTPError(http.StatusBadRequest, "days must be a number")
		}
	}

	res, err := api.core.MiningBlockWeightXFeeByHour(days)
	if err != nil {
		api.logger.Error("Error on core.MiningBlockWeightXFeeByHour", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &res)
}

func (api *API) MiningRevenueByHourLast(c echo.Context) error {
	var days int
	var err error
	daysStr := c.QueryParam("days")

	if daysStr != "" {
		days, err = strconv.Atoi(daysStr)
		if err != nil {
			api.logger.Error("days must be a number", zap.Error(err), zap.Any("days", daysStr))
			return echo.NewHTTPError(http.StatusBadRequest, "days must be a number")
		}
	}

	res, err := api.core.MiningRevenueByHour(days)
	if err != nil {
		api.logger.Error("Error on core.MiningRevenueByHour", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &res)
}
