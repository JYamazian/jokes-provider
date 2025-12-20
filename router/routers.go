package routes

import (
	"jokes-provider/controllers"
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all routes with controllers
func RegisterRoutes(app *fiber.App) {
	// Initialize controllers
	jokeCtrl := controllers.NewJokeController()
	healthCtrl := controllers.NewHealthController()
	metadataCtrl := controllers.NewMetadataController()

	// API v1 group
	v1 := app.Group("/api/v1")
	{
		// Jokes group
		jokes := v1.Group("/jokes")
		{
			jokes.Get("/random", jokeCtrl.GetRandomJoke)
		}

		// Metadata group
		v1.Get("/metadata", metadataCtrl.GetMetadata)
	}

	// Health group (outside API versioning)
	health := app.Group("/health")
	{
		health.Get("/readiness", healthCtrl.Readiness)
		health.Use(healthCtrl.SetupLivenessProbe("/liveness"))
	}

	// Swagger
	services.SetupSwagger(app)
}
