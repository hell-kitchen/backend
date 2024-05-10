package http

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strings"
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
		data     []string
		token    string
		ok       bool
		err      error
		userData *model.UserDataInToken
	)
	data, ok = c.Request().Header["Authorization"]
	if !ok {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}
	ctrl.log.Info("got authorization header", zap.Any("header_data", data))
	if len(data) == 0 {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}
	token = strings.Split(data[0], " ")[1]

	userData, err = ctrl.token.GetDataFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting access token: %w", err)
	}
	return userData, nil
}

func (ctrl *Controller) getUserIDFromAccessToken(c echo.Context) (uuid.UUID, error) {
	userData, err := ctrl.getUserDataFromRequest(c)
	if err != nil {
		return uuid.Nil, errors.New("unauthorized")
	}
	if !userData.IsAccess {
		return uuid.Nil, errors.New("unauthorized")
	}
	return userData.ID, nil
}
