package http

import (
	"github.com/hell-kitchen/backend/internal/model"
	"github.com/labstack/echo/v4"
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
	panic("not implemented")
}

func (ctrl *Controller) HandleRegister(c echo.Context) error {
	var request model.UsersLoginRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return nil
}

func (ctrl *Controller) HandleLogin(c echo.Context) error {
	panic("not implemented")
}

func (ctrl *Controller) HandleGetMe(c echo.Context) error {
	panic("not implemented")
}
