package http

import (
	"fmt"
	"github.com/google/uuid"
)

func (ctrl *Controller) generateAccessAndRefreshForUser(user uuid.UUID) (access string, refresh string, err error) {
	access, err = ctrl.token.CreateTokenForUser(user, true)
	if err != nil {
		return "", "", fmt.Errorf("error while creating access token: %w", err)
	}
	refresh, err = ctrl.token.CreateTokenForUser(user, false)
	if err != nil {
		return "", "", fmt.Errorf("error while creating refresh token: %w", err)
	}
	return
}
