package contoller

import "context"

type Controller interface {
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
}
