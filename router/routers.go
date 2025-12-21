package routes

import (
	"jokes-provider/controllers"
	"jokes-provider/services"
	"jokes-provider/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all routes with controllers
func RegisterRoutes(app *fiber.App) {
	// Initialize controllers
	jokeCtrl := controllers.NewJokeController()
	healthCtrl := controllers.NewHealthController()
	metadataCtrl := controllers.NewMetadataController()

	// API v1 group
	v1 := app.Group(utils.APIVersionV1)
	{
		// Jokes group
		jokes := v1.Group(utils.RouteJokes)
		{
			jokes.Get(utils.RandomJokeEndpoint, jokeCtrl.GetRandomJoke)
			jokes.Get(utils.JokeByIDEndpoint, jokeCtrl.GetJokeByID)
		}

		// Metadata group
		v1.Get(utils.MetadataEndpoint, metadataCtrl.GetMetadata)
	}

	// Health group (outside API versioning)
	health := app.Group(utils.RouteHealth)
	{
		health.Get(utils.ReadinessEndpoint, healthCtrl.Readiness)
		health.Use(healthCtrl.SetupLivenessProbe(utils.LivenessEndpoint))
	}

	// Swagger
	services.SetupSwagger(app)
}
