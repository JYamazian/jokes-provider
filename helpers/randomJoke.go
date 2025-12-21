package helpers

import (
	"encoding/csv"
	"errors"
	"jokes-provider/config"
	"jokes-provider/utils"
	"math/rand"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ErrJokeNotFound is returned when a joke with the specified ID is not found
var ErrJokeNotFound = errors.New("joke not found")

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

// GetJokeByID retrieves a joke by its ID from the CSV file
func GetJokeByID(c *fiber.Ctx, jokeID string) (map[string]string, error) {
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
		return nil, ErrJokeNotFound
	}

	// Get headers from first row
	headers := records[0]

	// Find the ID column index
	idIndex := -1
	for i, header := range headers {
		if header == utils.CSVColumnID {
			idIndex = i
			break
		}
	}

	if idIndex == -1 {
		config.LogError(c, "ID column not found in CSV file")
		return nil, errors.New(utils.ErrMsgIDColumnNotFound)
	}

	// Search for the joke with matching ID
	for _, row := range records[1:] {
		if idIndex < len(row) && row[idIndex] == jokeID {
			joke := make(map[string]string)
			for i, header := range headers {
				if i < len(row) {
					joke[header] = row[i]
				}
			}
			return joke, nil
		}
	}

	config.LogInfo(c, "Joke not found", "id", jokeID)
	return nil, ErrJokeNotFound
}
