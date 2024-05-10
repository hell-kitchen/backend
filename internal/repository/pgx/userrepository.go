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

func (repo *userRepository) migrate() error {
	query := `CREATE TABLE IF NOT EXISTS users
(
    id                 UUID PRIMARY KEY NOT NULL UNIQUE,
    username           VARCHAR          NOT NULL UNIQUE,
    encrypted_password VARCHAR          NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON users (username);`
	_, err := repo.pool.Exec(context.Background(), query)
	return err
}

func (repo *userRepository) Create(ctx context.Context, user *model.UserDTO) error {
	const query = `INSERT INTO users (id, username, encrypted_password) VALUES ($1, $2, $3);`
	_, err := repo.pool.Exec(ctx, query, user.ID, user.Username, user.Password)
	return err
}

func (repo *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.UserDTO, error) {
	u := new(model.UserDTO)
	const query = `SELECT u.id, u.username, u.encrypted_password
FROM users AS u
WHERE u.id = $1;`

	err := repo.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (*model.UserDTO, error) {
	u := new(model.UserDTO)
	const query = `SELECT u.id, u.username, u.encrypted_password
FROM users AS u
WHERE u.username = $1;`

	err := repo.pool.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}

	return u, nil
}
