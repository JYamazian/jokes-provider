package services

import (
	"jokes-provider/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MetadataHandler returns application metadata and configuration
func MetadataHandler(c *fiber.Ctx) error {
	config.LogInfo(c, "Metadata requested")

	metadata := fiber.Map{
		"app": fiber.Map{
			"name":    "Jokes Provider API",
			"version": config.AppConfig.Version,
			"flavor":  config.AppConfig.Flavor,
		},
		"server": fiber.Map{
			"port":        config.AppConfig.Port,
			"environment": config.AppConfig.Environment,
			"timestamp":   time.Now().Format(time.RFC3339),
		},
		"logging": fiber.Map{
			"level":          config.AppConfig.LogLevel,
			"format":         config.AppConfig.LogFormat,
			"format_type":    config.AppConfig.LogFormatType,
			"disable_colors": config.AppConfig.LogDisableColors,
		},
		"cache": fiber.Map{
			"enabled": config.CacheConfig.CacheEnabled,
			"url":     config.CacheConfig.CacheURL,
			"ttl":     config.CacheConfig.CacheTTL,
		},
		"files": fiber.Map{
			"jokes_path": config.AppConfig.JokesFilePath,
		},
		"headers": fiber.Map{
			"ip_header_name":      config.AppConfig.IPHeaderName,
			"country_header_name": config.AppConfig.CountryHeaderName,
		},
		"rate_limiter": fiber.Map{
			"max_requests": config.AppConfig.RateLimitMaxRequests,
			"duration":     config.AppConfig.RateLimitDuration,
		},
		"fiber": fiber.Map{
			"prefork":        config.AppConfig.FiberConfig.Prefork,
			"case_sensitive": config.AppConfig.FiberConfig.CaseSensitive,
			"strict_routing": config.AppConfig.FiberConfig.StrictRouting,
		},
	}

	return c.Status(fiber.StatusOK).JSON(metadata)
}
