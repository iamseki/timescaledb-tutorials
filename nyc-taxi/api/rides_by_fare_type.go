package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) ridesByFareType(c echo.Context) error {
	date := c.QueryParam("date")

	if date == "" {
		api.logger.Warn(fmt.Sprintf("bad request for date: '%v' on route %v", date, c.Request().RequestURI))
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := api.repository.RidesByFareType(date)
	if err != nil {
		api.logger.Fatal("error on api.repository.RidesByFareType", zap.Error(err))
	}

	return c.JSON(http.StatusOK, &res)
}
