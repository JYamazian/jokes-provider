package routes

import (
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// JokeRoute registers all joke-related routes
func JokeRoute(app *fiber.App) {
	app.Get("/joke/random", services.GetRandomJokeHandler)
}

// HealthRoute registers health-related routes (readiness only, liveness handled by middleware)
func HealthRoute(app *fiber.App) {
	app.Get("/health/readiness", services.ReadinessHandler)
}

// MetadataRoute registers metadata endpoint
func MetadataRoute(app *fiber.App) {
	app.Get("/metadata", services.MetadataHandler)
}

// SwaggerRoute sets up swagger middleware
func SwaggerRoute(app *fiber.App) {
	services.SetupSwagger(app)
}
