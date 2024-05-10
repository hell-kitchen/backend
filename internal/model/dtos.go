package model

import "github.com/google/uuid"

type UserDTO struct {
	ID       uuid.UUID
	Username string
	Password string
}

type TodoDTO struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsCompleted bool
	CreatedBy   uuid.UUID
}
