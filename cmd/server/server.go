package main

import (
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()
}
