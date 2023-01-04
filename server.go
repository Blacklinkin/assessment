package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
