package main

import (
	"context"
	"github.com/hell-kitchen/backend/internal/contoller"
	"github.com/hell-kitchen/backend/internal/contoller/http"
	"go.uber.org/zap"
)

func main() {
	var (
		ctx    context.Context
		log    *zap.Logger
		err    error
		server contoller.Controller
	)
	ctx = context.Background()

	log, _ = zap.NewProduction()

	server, err = http.New(log)
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
