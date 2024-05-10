package postgres

import (
	"context"
	"fmt"
	pgxzap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"go.uber.org/zap"
)

// New opens new postgres connection, configures it and return prepared client.
func New(connString string, log *zap.Logger) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	log.Info("initializing postgres client")

	c, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error while parsing db uri: %w", err)
	}

	var lvl = tracelog.LogLevelError
	c.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxzap.NewLogger(log),
		LogLevel: lvl,
	}
	c.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("postgres: init pgxpool: %w", err)
	}

	log.Info("created postgres client")
	return pool, nil
}
