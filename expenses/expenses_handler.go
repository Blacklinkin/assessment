package expenses

import (
	"net/http"
	"strconv"

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
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, h.Database.insertExpenses(*exp))
}

func (h *Handler) ViewExpensesByID(c echo.Context) error {
	if id := c.Param("id"); id != "" {
		Id, _ := strconv.Atoi(id)
		return c.JSON(http.StatusBadRequest, h.Database.viewExpensesDataByID(Id))
	}
	return c.JSON(http.StatusBadRequest, "invalid or forgot insert id")
}
