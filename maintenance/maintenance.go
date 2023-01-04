package maintenance

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
