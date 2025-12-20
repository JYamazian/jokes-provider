package api

import (
	"fmt"
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/middleware"
	routes "jokes-provider/router"
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// Initialize sets up and returns a configured Fiber application
func Initialize() (*fiber.App, error) {
	config.LoadEnvVars()

	app := fiber.New(fiber.Config{
		Prefork:       config.FiberConfig.Prefork,
		CaseSensitive: config.FiberConfig.CaseSensitive,
		StrictRouting: config.FiberConfig.StrictRouting,
		ServerHeader:  "Go Fiber - Jokes Provider",
		AppName:       "Jokes Provider API",
	})

	config.InitializeLogger(app)
	config.LogStartupInfo(config.AppConfig.Version, config.AppConfig.Flavor)

	if err := initRedis(); err != nil {
		return nil, err
	}

	if err := initJokesData(); err != nil {
		return nil, err
	}

	initMiddleware(app)
	routes.RegisterRoutes(app)

	return app, nil
}

// initRedis initializes the Redis connection
func initRedis() error {
	if err := middleware.InitRedis(); err != nil {
		config.LogError(nil, "Failed to initialize Redis connection", "error", err.Error())
		return fmt.Errorf("redis initialization failed: %w", err)
	}
	return nil
}

// initJokesData loads jokes from CSV file
func initJokesData() error {
	if err := helpers.LoadJokesFromCSV(nil, config.AppConfig.JokesFilePath); err != nil {
		config.LogError(nil, "Failed to load jokes data", "error", err.Error())
		return fmt.Errorf("jokes data loading failed: %w", err)
	}
	return nil
}

// initMiddleware sets up all middleware
func initMiddleware(app *fiber.App) {
	app.Use(services.SetupRateLimiter())
}

// Start starts the Fiber application server
func Start(app *fiber.App) error {
	return app.Listen(":" + config.AppConfig.Port)
}

// Shutdown gracefully shuts down the application
func Shutdown() error {
	if err := middleware.CloseRedis(); err != nil {
		config.LogError(nil, "Error closing Redis", "error", err.Error())
		return err
	}
	return nil
}
