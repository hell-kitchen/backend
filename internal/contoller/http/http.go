package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/config"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/hell-kitchen/backend/internal/pkg/token"
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"time"
)

var _ contoller.Controller = (*Controller)(nil)

type Controller struct {
	server *echo.Echo
	log    *zap.Logger
	cfg    *config.Config
	token  token.ProviderI
	repo   repository.Interface
}

func New(
	log *zap.Logger,
	cfg *config.Config,
	tokenProvider token.ProviderI,
	repo repository.Interface,
) (*Controller, error) {
	ctrl := &Controller{
		server: echo.New(),
		log:    log,
		cfg:    cfg,
		token:  tokenProvider,
		repo:   repo,
	}

	ctrl.configure()

	return ctrl, nil
}

func (ctrl *Controller) configureRoutes() {
	ctrl.server.GET("/ping", ctrl.HandlePing)

	api := ctrl.server.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", ctrl.HandleRegister)
			users.POST("/login", ctrl.HandleLogin)
			users.GET("/me", ctrl.HandleGetMe)
		}
		todos := api.Group("/todos")
		{
			todos.GET("/", ctrl.HandleGetTodos)
			todos.GET("", ctrl.HandleGetTodos)
			todos.POST("/", ctrl.HandleCreateTodo)
			todos.POST("", ctrl.HandleCreateTodo)
		}
	}
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
	ctrl.server.HideBanner = true
	ctrl.configureMiddlewares()
	ctrl.configureRoutes()
}

func (ctrl *Controller) OnStart(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		ctrl.log.Info("starting HTTP server", zap.String("bind-address", ctrl.cfg.Controller.GetBindAddress()))
		err := ctrl.server.Start(ctrl.cfg.Controller.GetBindAddress())
		if err != nil {
			cancel()
		}
	}()
	time.Sleep(300 * time.Millisecond)

	return ctx.Err()

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
