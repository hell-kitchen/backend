package config

import (
	"context"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
	"os"
	"path"
)

type Config struct {
	JWT        *JWT
	Controller *Controller
}

func New(ctx context.Context) (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		JWT: &JWT{
			AccessTokenLifetime:  20,
			RefreshTokenLifetime: 10000,
			PublicKeyPath:        path.Join(wd, "certs", "public.pem"),
			PrivateKeyPath:       path.Join(wd, "certs", "private.pem"),
		},
		Controller: &Controller{
			BindAddress: "0.0.0.0",
			BindPort:    8000,
		},
	}
	loader := confita.NewLoader(env.NewBackend(), flags.NewBackend())
	if err := loader.Load(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
