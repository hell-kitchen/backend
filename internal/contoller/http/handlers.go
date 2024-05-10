package http

import (
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (ctrl *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (ctrl *Controller) HandleGetTodos(c echo.Context) error {
	panic("not implemented")
}

func (ctrl *Controller) HandleGetTodoByID(c echo.Context) error {
	panic("not implemented")
}

func (ctrl *Controller) HandleCreateTodo(c echo.Context) error {
	panic("not implemented")
}

func (ctrl *Controller) HandleRegister(c echo.Context) error {
	var request model.UsersLoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	generated, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	user := &model.UserDTO{
		ID:       uuid.New(),
		Username: request.Username,
		Password: string(generated),
	}

	if err = ctrl.repo.User().Create(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshForUser(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	cookie := &http.Cookie{
		Name:    ctrl.cfg.AccessCookieName,
		Value:   accessToken,
		Expires: time.Now().Add(time.Duration(ctrl.cfg.AccessCookieLifetime) * time.Minute),
	}

	c.SetCookie(cookie)

	response := model.UsersLoginResponse{
		ID:           user.ID,
		RefreshToken: refreshToken,
	}

	return c.JSON(http.StatusOK, response)
}

func (ctrl *Controller) HandleLogin(c echo.Context) error {
	panic("not implemented")
}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	panic("not implemented")
}
