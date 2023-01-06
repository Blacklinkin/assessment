package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Database database
}

func (h *Handler) InitialDB() {
	h.Database.InitDatabase()
}

func (h *Handler) CloseDB() {
	h.Database.CloseDatabase()
}

func (h *Handler) AddExpenses(c echo.Context) error {
	exp := new(Expenses)
	if err := c.Bind(exp); err != nil {
		ErrMsg := Err{Message: err.Error()}
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, ErrMsg.Message)
	}
	result := h.Database.insertExpenses(*exp)
	return c.JSON(http.StatusCreated, result)
}

func (h *Handler) ViewExpensesByID(c echo.Context) error {
	if id := c.Param("id"); id != "" {
		return nil
	}
	return nil
}
