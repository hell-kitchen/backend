package pgx

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type userRepository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

var _ repository.UserRepository = (*userRepository)(nil)

func newUserRepository(pool *pgxpool.Pool, log *zap.Logger) (*userRepository, error) {
	repo := &userRepository{
		pool: pool,
		log:  log,
	}
	if err := repo.migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (u *userRepository) migrate() error {
	query := `CREATE TABLE IF NOT EXISTS users
(
    id                 UUID PRIMARY KEY NOT NULL UNIQUE,
    username           VARCHAR          NOT NULL UNIQUE,
    encrypted_password VARCHAR          NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON users (username);`
	_, err := u.pool.Exec(context.Background(), query)
	return err
}

func (u *userRepository) Create(ctx context.Context, user *model.UserDTO) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepository) GetByUsername(ctx context.Context, id uuid.UUID) (*model.UserDTO, error) {
	//TODO implement me
	panic("implement me")
}
