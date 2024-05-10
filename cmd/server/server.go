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

	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	server, err = http.New()
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Error(
			"stopped server",
			zap.Error(server.OnStop(ctx)),
		)
	}()
	err = server.OnStart(ctx)
	if err != nil {
		log.Error("got error on server startup", zap.Error(err))
	}
}
