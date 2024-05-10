package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var _ contoller.Controller = (*Controller)(nil)

type Controller struct {
	server *echo.Echo
	log    *zap.Logger
}

func New(log *zap.Logger) (*Controller, error) {
	ctrl := &Controller{
		server: echo.New(),
		log:    log,
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
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogValuesFunc: ctrl.logValuesFunc,
			LogLatency:    true,
			LogRequestID:  true,
			LogMethod:     true,
			LogURI:        true,
		}),
	}
	ctrl.server.Use(middlewares...)
}

func (ctrl *Controller) configure() {
	ctrl.configureMiddlewares()
	ctrl.configureRoutes()
}

func (ctrl *Controller) OnStart(_ context.Context) error {
	ctrl.log.Info("starting HTTP server on port :8080")
	return ctrl.server.Start(":8080")
}

func (ctrl *Controller) OnStop(ctx context.Context) error {
	return ctrl.server.Shutdown(ctx)
}

func (ctrl *Controller) logValuesFunc(_ echo.Context, v middleware.RequestLoggerValues) error {
	ctrl.log.Info(
		"Request",
		zap.String("uri", v.URI),
		zap.String("method", v.Method),
		zap.Duration("duration", v.Latency),
		zap.String("request-id", v.RequestID),
	)
	return nil
}
