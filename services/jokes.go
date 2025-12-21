package services

import (
	"jokes-provider/helpers"
	"jokes-provider/utils"
	"jokes-provider/wrapper"

	"github.com/gofiber/fiber/v2"
)

type JokeService struct{}

func NewJokeService() *JokeService {
	return &JokeService{}
}

func (s *JokeService) GetRandomJoke(c *fiber.Ctx, cacheKey string) (map[string]string, error) {
	// Try cache first
	if cached, ok := wrapper.ReadCacheIfAllowed(c, cacheKey); ok {
		return cached, nil
	}

	// Cache miss - get from data source
	joke, err := helpers.GetRandomJoke(c)
	if err != nil {
		return nil, err
	}

	// Write to cache
	_ = wrapper.WriteCacheIfAllowed(c, cacheKey, joke)

	return joke, nil
}

func (s *JokeService) GetJokeByID(c *fiber.Ctx, jokeID string) (map[string]string, error) {
	cacheKey := utils.CacheKeyPrefixJoke + jokeID

	if cached, ok := wrapper.ReadCacheIfAllowed(c, cacheKey); ok {
		return cached, nil
	}

	joke, err := helpers.GetJokeByID(c, jokeID)
	if err != nil {
		return nil, err
	}

	_ = wrapper.WriteCacheIfAllowed(c, cacheKey, joke)

	return joke, nil
}
