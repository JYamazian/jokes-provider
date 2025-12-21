package services

import (
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

// HealthService handles health check business logic
type HealthService struct{}

// NewHealthService creates a new HealthService instance
func NewHealthService() *HealthService {
	return &HealthService{}
}

// CheckReadiness checks if the service is ready
func (s *HealthService) CheckReadiness(c *fiber.Ctx) models.ReadinessHealthStatus {
	// Check Redis status
	if !helpers.CheckRedisStatus(c) {
		config.LogError(c, "Readiness check failed: Redis unavailable")
		return models.ReadinessHealthStatus{
			Ready:  false,
			Reason: "Redis unavailable",
		}
	}

	// Check if CSV file is accessible
	if !config.FileExists(config.AppConfig.JokesFilePath) {
		config.LogError(c, "Readiness check failed: Jokes CSV file not accessible", "path", config.AppConfig.JokesFilePath)
		return models.ReadinessHealthStatus{
			Ready:  false,
			Reason: "Jokes CSV file not accessible",
		}
	}

	config.LogInfo(c, "Readiness check passed")
	return models.ReadinessHealthStatus{
		Ready: true,
		Redis: "connected",
		CSV:   "accessible",
	}
}

// SetupLivenessProbe returns Fiber's built-in healthcheck middleware for liveness
func (s *HealthService) SetupLivenessProbe(endpoint string) fiber.Handler {
	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: endpoint,
	})
}
