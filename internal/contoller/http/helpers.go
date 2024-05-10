package http

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
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

func (ctrl *Controller) getUserDataFromRequest(c echo.Context) (*model.UserDataInToken, error) {
	var (
		cookie   *http.Cookie
		err      error
		userData *model.UserDataInToken
	)
	cookie, err = c.Cookie(ctrl.cfg.AccessCookieName)
	if err != nil {
		return nil, fmt.Errorf("error while getting access cookie: %w", err)
	}

	userData, err = ctrl.token.GetDataFromToken(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}
	return userData, nil
}
