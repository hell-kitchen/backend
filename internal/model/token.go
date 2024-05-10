package model

import "github.com/google/uuid"

type UserDataInToken struct {
	ID       uuid.UUID
	IsAccess bool
}
