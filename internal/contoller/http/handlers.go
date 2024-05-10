package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (ctrl *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
