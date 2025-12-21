package controllers

import (
	"jokes-provider/config"
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// HealthController handles health check endpoints
type HealthController struct {
	healthService *services.HealthService
}

// NewHealthController creates a new HealthController instance
func NewHealthController() *HealthController {
	return &HealthController{
		healthService: services.NewHealthService(),
	}
}

// Readiness godoc
// @Summary      Readiness check
// @Description  Checks if the service is ready and dependencies are available (Redis, CSV file)
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ReadinessHealthStatus  "Service is ready"
// @Failure      503  {object}  models.ReadinessHealthStatus  "Service is not ready"
// @Router       /health/readiness [get]
func (ctrl *HealthController) Readiness(c *fiber.Ctx) error {
	config.LogInfo(c, "Readiness check called")

	status := ctrl.healthService.CheckReadiness(c)

	if !status.Ready {
		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	}

	return c.Status(fiber.StatusOK).JSON(status)
}

// SetupLivenessProbe returns Fiber's built-in liveness probe middleware
// @Summary      Liveness probe
// @Description  Simple liveness probe to check if the service is running (uses Fiber built-in)
// @Tags         health
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /health/liveness [get]
func (ctrl *HealthController) SetupLivenessProbe(endpoint string) fiber.Handler {
	return ctrl.healthService.SetupLivenessProbe(endpoint)
}
