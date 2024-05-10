package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ contoller.Controller = (*Controller)(nil)

type Controller struct {
	server *echo.Echo
}

func New() (*Controller, error) {
	ctrl := &Controller{
		server: echo.New(),
	}

	ctrl.configure()

	return ctrl, nil
}

func (ctrl *Controller) configureRoutes() {
	ctrl.server.GET("/ping", ctrl.HandlePing)
}

func (ctrl *Controller) configureMiddlewares() {
	var middlewares = []echo.MiddlewareFunc{
		middleware.Gzip(),
		middleware.Recover(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Skipper:      middleware.DefaultSkipper,
			Generator:    uuid.NewString,
			TargetHeader: echo.HeaderXRequestID,
		}),
		middleware.Logger(),
	}
	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) configure() {
	ctrl.configureMiddlewares()
	ctrl.configureRoutes()
}

func (ctrl *Controller) OnStart(_ context.Context) error {
	return ctrl.server.Start(":8080")
}

func (ctrl *Controller) OnStop(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}
