package services

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func SetupSwagger(app *fiber.App) {
	// Setup Swagger documentation as middleware
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Jokes Provider API",
	}))
}
