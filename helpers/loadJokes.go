package helpers

import (
	"jokes-provider/config"

	"github.com/gofiber/fiber/v2"
)

// LoadJokesFromCSV validates that the CSV file is accessible (no longer caches in memory)
func LoadJokesFromCSV(c *fiber.Ctx, filePath string) error {
	if !config.FileExists(filePath) {
		config.LogError(c, "CSV file not found", "file_path", filePath)
		return nil
	}

	data, err := config.ReadCSVWithHeaders(c, filePath)
	if err != nil {
		config.LogError(c, "Failed to validate CSV file", "file_path", filePath, "error", err.Error())
		return nil
	}

	if len(data) == 0 {
		config.LogError(c, "CSV file is empty", "file_path", filePath)
		return nil
	}

	config.LogInfo(c, "CSV file validated", "file_path", filePath, "joke_count", len(data))
	return nil
}
