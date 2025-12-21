package controllers

import (
	"jokes-provider/config"
	"jokes-provider/services"

	"github.com/gofiber/fiber/v2"
)

// MetadataController handles metadata endpoints
type MetadataController struct {
	metadataService *services.MetadataService
}

// NewMetadataController creates a new MetadataController instance
func NewMetadataController() *MetadataController {
	return &MetadataController{
		metadataService: services.NewMetadataService(),
	}
}

// GetMetadata godoc
// @Summary      Get application metadata
// @Description  Returns comprehensive application metadata including version, configuration, and environment information
// @Tags         metadata
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Metadata  "Application metadata"
// @Router       /api/v1/metadata [get]
func (ctrl *MetadataController) GetMetadata(c *fiber.Ctx) error {
	config.LogInfo(c, "Metadata requested")

	metadata := ctrl.metadataService.GetMetadata()

	return c.Status(fiber.StatusOK).JSON(metadata)
}
