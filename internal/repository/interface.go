package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.UserDTO) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserDTO, error)
	GetByUsername(ctx context.Context, username string) (*model.UserDTO, error)
}

type TodosRepository interface {
	Create(ctx context.Context, todo *model.TodoDTO) error
	GetAll(ctx context.Context) ([]model.TodoDTO, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.TodoDTO, error)
}

type Interface interface {
	User() UserRepository
	Todos() TodosRepository
}
