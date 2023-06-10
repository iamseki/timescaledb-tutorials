package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (api *API) ridesByInterval(c echo.Context) error {
	date := c.QueryParam("date")
	interval := c.QueryParam("interval")

	if date == "" || interval == "" {
		api.logger.Warn(fmt.Sprintf("bad request on route %v", c.Request().RequestURI))
		return c.NoContent(http.StatusBadRequest)
	}

	res, err := api.repository.RidesByTimeBucket(date, interval)
	if err != nil {
		api.logger.Fatal("error on api.repository.RidesByTimeBucket", zap.Error(err))
	}

	return c.JSON(http.StatusOK, &res)
}
