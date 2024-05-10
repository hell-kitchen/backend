package main

import (
	"context"
	"github.com/hell-kitchen/backend/internal/config"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/hell-kitchen/backend/internal/contoller/http"
	"github.com/hell-kitchen/backend/internal/pkg/token/jwt"
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
	)
	ctx = context.Background()

	log, _ = zap.NewProduction()

	cfg, err = config.New(ctx)
	if err != nil {
		log.Fatal("Failed to create config", zap.Error(err))
	}
	log.Info("xd", zap.Any("cfg", cfg))

	provider, err = jwt.NewProvider(cfg.JWT, log)
	if err != nil {
		log.Fatal("Failed to create jwt provider", zap.Error(err))
	}

	server, err = http.New(log, cfg.Controller, provider)
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
