package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) ridesByAirportsCode(c echo.Context) error {
	date := c.QueryParam("date")
	airportsCode := c.QueryParam("airportsCode")

	if date == "" {
		api.logger.Warn(fmt.Sprintf("bad request for date: '%v' on route %v", date, c.Request().RequestURI))
		return c.NoContent(http.StatusBadRequest)
	}

	if airportsCode == "" {
		api.logger.Warn(fmt.Sprintf("bad request for airportsCode: '%v' on route %v", airportsCode, c.Request().RequestURI))
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := api.repository.RidesByAirportsCodeResponse(date, airportsCode)
	if err != nil {
		api.logger.Fatal("error on api.repository.AverageFareSince", zap.Error(err))
	}

	return c.JSON(http.StatusOK, &res)
}
