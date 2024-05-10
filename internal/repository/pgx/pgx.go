package pgx

import (
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var _ repository.Interface = (*Repository)(nil)

type Repository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
	todo *todoRepository
	user *userRepository
}

func New(pool *pgxpool.Pool, log *zap.Logger) (*Repository, error) {
	users, err := newUserRepository(pool, log)
	if err != nil {
		return nil, err
	}

	todos, err := newTodoRepository(pool, log)
	if err != nil {
		return nil, err
	}

	repo := &Repository{
		pool: pool,
		log:  log,
		todo: todos,
		user: users,
	}
	return repo, nil
}

func (r *Repository) User() repository.UserRepository {
	return r.user
}

func (r *Repository) Todos() repository.TodosRepository {
	return r.todo
}
