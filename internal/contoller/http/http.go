package http

import (
	"context"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/labstack/echo/v4"
)

var _ contoller.Controller = (*Controller)(nil)

type Controller struct {
	server *echo.Echo
}

func New() (*Controller, error) {
	ctrl := &Controller{
		server: echo.New(),
	}
	return ctrl, nil
}

func (ctrl *Controller) OnStart(_ context.Context) error {
	return ctrl.server.Start(":8080")
}

func (ctrl *Controller) OnStop(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
