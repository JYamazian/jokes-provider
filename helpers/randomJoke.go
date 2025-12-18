package helpers

import (
	"encoding/csv"
	"jokes-provider/config"
	"math/rand"
	"os"

	"github.com/gofiber/fiber/v2"
)

// GetRandomJoke reads a random joke directly from the CSV file
func GetRandomJoke(c *fiber.Ctx) (map[string]string, error) {
	file, err := os.Open(config.AppConfig.JokesFilePath)
	if err != nil {
		config.LogError(c, "Error opening jokes CSV file", "path", config.AppConfig.JokesFilePath, "error", err.Error())
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		config.LogError(c, "Error reading CSV file", "error", err.Error())
		return nil, err
	}

	// Need at least header + 1 joke
	if len(records) < 2 {
		config.LogError(c, "No jokes available in CSV file")
		return nil, nil
	}

	// Get headers from first row
	headers := records[0]

	// Pick random joke from index 1 onwards (skip header)
	randomIndex := rand.Intn(len(records)-1) + 1
	jokeRow := records[randomIndex]

	// Build joke map from headers and values
	joke := make(map[string]string)
	for i, header := range headers {
		if i < len(jokeRow) {
			joke[header] = jokeRow[i]
		}
	}

	return joke, nil
}
