package pgx

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/hell-kitchen/backend/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type todoRepository struct {
	pool *pgxpool.Pool
	log  *zap.Logger
}

func newTodoRepository(pool *pgxpool.Pool, log *zap.Logger) (*todoRepository, error) {
	repo := &todoRepository{
		pool: pool,
		log:  log,
	}
	if err := repo.migrate(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (repo *todoRepository) migrate() error {
	query := `CREATE TABLE IF NOT EXISTS todos
(
    "id"           UUID PRIMARY KEY NOT NULL UNIQUE,
    "name"         VARCHAR          NOT NULL UNIQUE,
    "description"  VARCHAR          NOT NULL,
    "is_completed" BOOLEAN          NOT NULL DEFAULT FALSE,
    "created_by"   UUID REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS todos_created_by_idx ON todos (created_by);`
	_, err := repo.pool.Exec(context.Background(), query)
	return err
}

func (repo *todoRepository) Create(ctx context.Context, todo *model.TodoDTO) error {
	const query = `INSERT INTO todos (id, name, description, created_by) VALUES ($1, $2, $3, $4);`
	_, err := repo.pool.Exec(ctx, query, todo.ID, todo.Name, todo.Description, todo.CreatedBy)
	return err
}

func (repo *todoRepository) GetByID(ctx context.Context) (*model.TodoDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *todoRepository) getAllByQuery(ctx context.Context, query string, args ...any) ([]model.TodoDTO, error) {
	var res []model.TodoDTO
	rows, err := repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp model.TodoDTO

		err = rows.Scan(&temp.ID, &temp.Name, &temp.Description, &temp.IsCompleted, &temp.CreatedBy)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *todoRepository) GetAll(ctx context.Context) ([]model.TodoDTO, error) {
	const query = `SELECT t.id, t.name, t.description, t.is_completed, t.created_by
FROM todos t;`
	return repo.getAllByQuery(ctx, query)
}

func (repo *todoRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.TodoDTO, error) {
	const query = `SELECT t.id, t.name, t.description, t.is_completed, t.created_by
FROM todos t
WHERE t.created_by = $1;`
	return repo.getAllByQuery(ctx, query, userID)
}

var _ repository.TodosRepository = (*todoRepository)(nil)
