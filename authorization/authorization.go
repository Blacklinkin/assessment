package authorization

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Process is the middleware function.
func AuthHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") != "November 10, 2009" {
			c.JSON(http.StatusUnauthorized, nil)
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
