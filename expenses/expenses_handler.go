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
	exp := Expenses{}
	if err := c.Bind(&exp); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, h.Database.insertExpenses(exp))
}

func (h *Handler) ViewExpensesByID(c echo.Context) error {
	if id := c.Param("id"); id != "" {
		Id, _ := strconv.Atoi(id)
		return c.JSON(http.StatusOK, h.Database.viewExpensesDataByID(Id))
	}
	return c.JSON(http.StatusBadRequest, "invalid or forgot insert id")
}

func (h *Handler) UpdateExpenses(c echo.Context) error {
	expUpdate := Expenses{}
	err := c.Bind(&expUpdate)
	if err != nil {
		return err
	}
	if id := c.Param("id"); id != "" {
		Id, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		h.Database.updateExpensesDataBase(Id, expUpdate)
		return c.JSON(http.StatusOK, h.Database.viewExpensesDataByID(Id))
	}
	return c.JSON(http.StatusBadRequest, err.Error())
}

func (h *Handler) ViewAllExpenses(c echo.Context) error {
	expSetResult := h.Database.viewAllExpenses()
	if expSetResult == nil {
		return c.JSON(http.StatusInternalServerError, h.Database.errMsg)
	}
	return c.JSON(http.StatusOK, expSetResult)
}
