package token

import (
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
)

type ProviderI interface {
	GetDataFromToken(token string) (*model.UserDataInToken, error)
	CreateTokenForUser(userID uuid.UUID, isAccess bool) (string, error)
}
