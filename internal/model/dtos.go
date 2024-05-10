package model

import "github.com/google/uuid"

type UserDTO struct {
	ID       uuid.UUID
	Username string
	Password string
}

type TodoDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	CreatedBy   uuid.UUID `json:"created_by"`
}
