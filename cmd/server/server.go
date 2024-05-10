package main

import (
	"context"
	"github.com/hell-kitchen/backend/internal/config"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/hell-kitchen/backend/internal/contoller/http"
	"github.com/hell-kitchen/backend/internal/pkg/token"
	"github.com/hell-kitchen/backend/internal/pkg/token/jwt"
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/hell-kitchen/backend/internal/repository/pgx"
	"github.com/hell-kitchen/backend/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fx.New(NewApp()).Run()
}

func NewApp() fx.Option {
	return fx.Options(
		fx.Provide(
			context.Background,
			zap.NewProduction,
			config.New,
			postgres.New,
			fx.Annotate(jwt.NewProvider, fx.As(new(token.ProviderI))),
			fx.Annotate(pgx.New, fx.As(new(repository.Interface))),
			fx.Annotate(http.New, fx.As(new(contoller.Controller))),
		),
		fx.Invoke(
			startHttp,
		),
	)
}

func startHttp(lc fx.Lifecycle, ctrl contoller.Controller) {
	lc.Append(fx.Hook{
		OnStart: ctrl.OnStart,
		OnStop:  ctrl.OnStop,
	})
}

func oldMain() {
	var (
		ctx      context.Context
		log      *zap.Logger
		err      error
		server   contoller.Controller
		provider *jwt.Provider
		cfg      *config.Config
		pool     *pgxpool.Pool
		repo     repository.Interface
	)
	ctx = context.Background()

	log, _ = zap.NewProduction()

	cfg, err = config.New(ctx)
	if err != nil {
		log.Fatal("Failed to create config", zap.Error(err))
	}
	log.Info("initialized config", zap.Any("cfg", cfg))

	provider, err = jwt.NewProvider(cfg, log)
	if err != nil {
		log.Fatal("Failed to create jwt provider", zap.Error(err))
	}

	pool, err = postgres.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to create postgres pool", zap.Error(err))
	}

	repo, err = pgx.New(pool, log)
	if err != nil {
		log.Fatal("error while creating pgx storage", zap.Error(err))
	}

	server, err = http.New(log, cfg, provider, repo)
	if err != nil {
		log.Fatal("Failed to create server", zap.Error(err))
	}
	defer func() {
		log.Error(
			"Stopping HTTP server",
			zap.Error(server.OnStop(ctx)),
		)
	}()
	err = server.OnStart(ctx)
	if err != nil {
		log.Fatal("got error on server startup", zap.Error(err))
	}

	var sig os.Signal
	interrupt := make(chan os.Signal, 1)
	closed := make(chan struct{})
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		sig = <-interrupt
		defer func() {
			close(closed)
		}()
		if err = server.OnStop(ctx); err != nil {
			log.Fatal("error while stopping http server", zap.Error(err))
		}
	}()
	<-closed
	log.Info("graceful shutdown", zap.String("signal", sig.String()))
}
