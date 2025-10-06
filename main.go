package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func startHeartbeat(e *echo.Echo) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			e.Logger.Info("Service heartbeat - sersver is running...")
		}
	}()
}

// checkMountPath verifies that a directory exists and is writable.
// It does this by creating and then deleting a temporary file.
func checkMountPath(path string, logger echo.Logger) {
	logger.Infof("Checking mount path at: %s", path)

	// 1. Check if the path exists and is a directory.
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		logger.Fatalf("Mount path does not exist: %s", path)
	}
	if !info.IsDir() {
		logger.Fatalf("Mount path is not a directory: %s", path)
	}

	// 2. Check if the path is writable by creating a temporary file.
	tmpFile := filepath.Join(path, ".tmp-write-check")
	if err := os.WriteFile(tmpFile, []byte("test"), 0600); err != nil {
		logger.Fatalf("Mount path is not writable: %s. Error: %v", path, err)
	}

	// 3. Clean up the temporary file.
	if err := os.Remove(tmpFile); err != nil {
		// This is not a fatal error, but we should log it.
		logger.Warnf("Could not remove temporary check file: %s. Error: %v", tmpFile, err)
	}

	logger.Infof("Successfully verified mount path is available and writable: %s", path)
}


func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.INFO)

	mountPath := os.Getenv("MOUNT_PATH")
	if mountPath == "" {
		e.Logger.Fatal("MOUNT_PATH environment variable not set. Exiting.")
	}
	// Run the check. This will call log.Fatal and exit if it fails.
	checkMountPath(mountPath, e.Logger)
	
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
