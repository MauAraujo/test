package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func startHeartbeat(e *echo.Echo) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			e.Logger.Info("Service heartbeat - sefrver is running...")
		}
	}()
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.INFO)

	// Start the heartbeat logging
	startHeartbeat(e)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker!! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	e.Logger.Info("Starting server on port " + httpPort)
	e.Logger.Fatal(e.Start(":" + httpPort))
}

// Simple implementation of an integer minimum
// Adapted from: https://gobyexample.com/testing-and-benchmarking
func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
