package main

import (
	"context"
	"github.com/hell-kitchen/backend/internal/config"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/hell-kitchen/backend/internal/contoller/http"
	"github.com/hell-kitchen/backend/internal/pkg/token/jwt"
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/hell-kitchen/backend/internal/repository/pgx"
	"github.com/hell-kitchen/backend/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
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

	provider, err = jwt.NewProvider(cfg.JWT, log)
	if err != nil {
		log.Fatal("Failed to create jwt provider", zap.Error(err))
	}

	pool, err = postgres.New(cfg.Postgres.DNS, log)
	if err != nil {
		log.Fatal("Failed to create postgres pool", zap.Error(err))
	}

	repo, err = pgx.New(pool, log)
	if err != nil {
		log.Fatal("error while creating pgx storage", zap.Error(err))
	}

	server, err = http.New(log, cfg.Controller, provider, repo)
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
		log.Error("got error on server startup", zap.Error(err))
	}
}
