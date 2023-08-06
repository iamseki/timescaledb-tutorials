package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthcheckResponse struct {
	OK bool `json:"ok"`
}

func (api *API) healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, &HealthcheckResponse{OK: true})
}
