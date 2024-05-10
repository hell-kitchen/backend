package model

import "github.com/google/uuid"

type UsersLoginResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}
