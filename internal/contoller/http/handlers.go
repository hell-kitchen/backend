package http

import (
	"github.com/google/uuid"
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
	user, err := ctrl.getUserIDFromAccessToken(c)
	if err != nil {
		ctrl.log.Error("could not validate access token from headers", zap.Error(err))
		return err
	}
	ctrl.log.Info("logged in", zap.String("user_id", user.String()))

	var request model.TodoCreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	todo := &model.TodoDTO{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		IsCompleted: false,
		CreatedBy:   user,
	}

	if err = ctrl.repo.Todos().Create(c.Request().Context(), todo); err != nil {
		ctrl.log.Error("error while creation todo", zap.Error(err))
		return err
	}
	ctrl.log.Info("successfully created todo", zap.Any("todo", todo))

	return c.JSON(http.StatusCreated, todo)
}

func (ctrl *Controller) HandleRegister(c echo.Context) error {
	var request model.UsersRegisterRequest
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	generated, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	user := &model.UserDTO{
		ID:       uuid.New(),
		Username: request.Username,
		Password: string(generated),
	}
	ctrl.log.Info("got user", zap.Any("user", user))

	if err = ctrl.repo.User().Create(c.Request().Context(), user); err != nil {
		ctrl.log.Error("got error while creating user", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.log.Info("successfully created user")

	accessToken, refreshToken, err := ctrl.generateAccessAndRefreshForUser(user.ID)
	if err != nil {
		ctrl.log.Error("error while creating tokens", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	response := model.UsersRegisterResponse{
		ID:           user.ID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) HandleLogin(c echo.Context) error {
	var request model.UsersLoginRequest
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("could not bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, err := ctrl.repo.User().GetByUsername(c.Request().Context(), request.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	access, refresh, err := ctrl.generateAccessAndRefreshForUser(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	response := &model.UsersLoginResponse{
		ID:           user.ID,
		AccessToken:  access,
		RefreshToken: refresh,
	}
	return c.JSON(http.StatusOK, response)
}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	panic("not implemented")
}
