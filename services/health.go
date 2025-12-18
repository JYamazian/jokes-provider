package services

import (
	"jokes-provider/config"
	"jokes-provider/helpers"

	"github.com/gofiber/fiber/v2"
)

// ReadinessHandler handles readiness check requests (checks dependencies)
func ReadinessHandler(c *fiber.Ctx) error {
	config.LogInfo(c, "Readiness check called")

	// Check Redis status
	if !helpers.CheckRedisStatus(c) {
		config.LogError(c, "Readiness check failed: Redis unavailable")
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"ready":  false,
			"reason": "Redis unavailable",
		})
	}

	// Check if CSV file is accessible
	if !config.FileExists(config.AppConfig.JokesFilePath) {
		config.LogError(c, "Readiness check failed: Jokes CSV file not accessible", "path", config.AppConfig.JokesFilePath)
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"ready":  false,
			"reason": "Jokes CSV file not accessible",
		})
	}

	config.LogInfo(c, "Readiness check passed")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ready": true,
		"redis": "connected",
		"csv":   "accessible",
	})
}
