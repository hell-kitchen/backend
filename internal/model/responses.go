package model

import "github.com/google/uuid"

type UsersLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	RefreshToken string    `json:"refresh_token"`
}
