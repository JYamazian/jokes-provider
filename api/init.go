package api

import (
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/middleware"
	routes "jokes-provider/router"
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// Initialize sets up and returns a configured Fiber application
func Initialize() *fiber.App {
	// Load environment variables
	config.LoadEnvVars()

	// Initialize Redis connection (singleton)
	if err := middleware.InitRedis(); err != nil {
		config.LogError(nil, "Failed to initialize Redis", "error", err.Error())
	}

	// Initialize Fiber app with config
	app := config.InitializeApp()

	// Log startup information with build details
	config.LogStartupInfo(config.AppConfig.Version, config.AppConfig.Flavor)

	// Load jokes from CSV
	if err := helpers.LoadJokesFromCSV(nil, config.AppConfig.JokesFilePath); err != nil {
		config.LogError(nil, "Warning: Could not load jokes from CSV", "error", err.Error())
	}

	// Setup rate limiter middleware globally
	app.Use(services.SetupRateLimiter())

	// Setup Swagger middleware first
	routes.SwaggerRoute(app)

	// Register all route groups from routes package
	routes.HealthRoute(app)
	routes.JokeRoute(app)
	routes.MetadataRoute(app)

	return app
}

// Start starts the Fiber application server
func Start(app *fiber.App) error {
	return app.Listen(":" + config.AppConfig.Port)
}
