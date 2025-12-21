package wrapper

import (
	"encoding/json"
	"jokes-provider/config"
	"jokes-provider/middleware"

	"github.com/gofiber/fiber/v2"
)

// WriteCacheIfAllowed writes data to cache if caching is enabled and allowed by headers
func WriteCacheIfAllowed(c *fiber.Ctx, cacheKey string, data map[string]string) error {
	if !config.CacheConfig.CacheEnabled {
		config.LogInfo(c, "Skipping cache WRITE - caching disabled", "cache_key", cacheKey)
		return nil
	}

	if shouldSkipCache(c) {
		config.LogInfo(c, "Skipping cache WRITE - Cache-Control: no-cache", "cache_key", cacheKey)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		config.LogError(c, "Error marshaling data for cache", "cache_key", cacheKey, "error", err.Error())
		return err
	}

	if err := middleware.SetToCache(c, cacheKey, jsonData); err != nil {
		return err
	}

	return nil
}

// ReadCacheIfAllowed reads data from cache if caching is enabled and allowed by headers
// Returns (data, cacheHit) - cacheHit is true if data was found in cache
func ReadCacheIfAllowed(c *fiber.Ctx, cacheKey string) (map[string]string, bool) {
	if !config.CacheConfig.CacheEnabled {
		config.LogInfo(c, "Skipping cache READ - caching disabled", "cache_key", cacheKey)
		return nil, false
	}

	if shouldSkipCache(c) {
		config.LogInfo(c, "Skipping cache READ - Cache-Control: no-cache", "cache_key", cacheKey)
		return nil, false
	}

	cachedData, err := middleware.GetFromCache(c, cacheKey)
	if err != nil {
		return nil, false
	}

	if cachedData == nil {
		return nil, false
	}

	var result map[string]string
	if err := json.Unmarshal(cachedData, &result); err != nil {
		config.LogError(c, "Error unmarshaling cached data", "cache_key", cacheKey, "error", err.Error())
		return nil, false
	}

	return result, true
}
