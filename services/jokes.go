package services

import (
	"encoding/json"
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/middleware"
	"path"

	"github.com/gofiber/fiber/v2"
)

// GetRandomJokeHandler handles requests for a random joke with caching
func GetRandomJokeHandler(c *fiber.Ctx) error {
	cacheKey := path.Base(c.Path())

	// Try to get from cache first using middleware if enabled
	if config.CacheConfig.CacheEnabled {
		cachedJoke, err := middleware.GetFromCache(c, cacheKey)
		if err == nil && cachedJoke != nil {
			// Cache hit - deserialize and return
			var joke map[string]string
			if err := json.Unmarshal(cachedJoke, &joke); err == nil {
				return c.Status(fiber.StatusOK).JSON(joke)
			}
		}
	}

	// Cache miss - get a new random joke
	joke, err := helpers.GetRandomJoke(c)
	if err != nil {
		config.LogError(c, "Error retrieving random joke", "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve joke",
		})
	}

	// Set in cache using middleware if enabled
	if config.CacheConfig.CacheEnabled {
		// Marshal joke to JSON for caching
		jokeJSON, err := json.Marshal(joke)
		if err != nil {
			config.LogError(c, "Error marshaling joke for cache", "error", err.Error())
			return c.Status(fiber.StatusOK).JSON(joke)
		}

		// Set in cache using middleware
		if err := middleware.SetToCache(c, cacheKey, jokeJSON); err != nil {
			// Cache failure is non-blocking - still return the joke
			config.LogInfo(c, "Returning joke despite cache error")
		}
	}

	return c.Status(fiber.StatusOK).JSON(joke)
}
