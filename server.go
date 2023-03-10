package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Blacklinkin/assessment/authorization"
	"github.com/Blacklinkin/assessment/expenses"
	"github.com/Blacklinkin/assessment/maintenance"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	//Setup
	h := expenses.Handler{}
	h.InitialDB()
	e := echo.New()

	//Add Middleware of e object
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(authorization.AuthHeader)

	//Initial Path
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hallo ,Expense Tracking System v1.0.0")
	})

	e.GET("/health", maintenance.HealthHandler)

	e.POST("/expenses", h.AddExpenses)
	e.GET("/expenses/:id", h.ViewExpensesByID)
	e.PUT("/expenses/:id", h.UpdateExpenses)
	e.GET("/expenses", h.ViewAllExpenses)
	//Graceful shutdown

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server!")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer h.CloseDB()
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("bye bye see you next time <(^^)>")
}
