package services

import (
	"encoding/json"
	"jokes-provider/config"
	"jokes-provider/helpers"
	"jokes-provider/middleware"

	"github.com/gofiber/fiber/v2"
)

// JokeService handles joke business logic
type JokeService struct{}

// NewJokeService creates a new JokeService instance
func NewJokeService() *JokeService {
	return &JokeService{}
}

// GetRandomJoke retrieves a random joke with caching support
func (s *JokeService) GetRandomJoke(c *fiber.Ctx, cacheKey string) (map[string]string, error) {
	// Try to get from cache first if enabled
	if config.CacheConfig.CacheEnabled {
		cachedJoke, err := middleware.GetFromCache(c, cacheKey)
		if err == nil && cachedJoke != nil {
			var joke map[string]string
			if err := json.Unmarshal(cachedJoke, &joke); err == nil {
				return joke, nil
			}
		}
	}

	// Cache miss - get a new random joke
	joke, err := helpers.GetRandomJoke(c)
	if err != nil {
		return nil, err
	}

	// Set in cache if enabled
	if config.CacheConfig.CacheEnabled {
		jokeJSON, err := json.Marshal(joke)
		if err != nil {
			config.LogError(c, "Error marshaling joke for cache", "error", err.Error())
			return joke, nil
		}

		if err := middleware.SetToCache(c, cacheKey, jokeJSON); err != nil {
			config.LogInfo(c, "Returning joke despite cache error")
		}
	}

	return joke, nil
}
