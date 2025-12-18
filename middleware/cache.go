package middleware

import (
	"jokes-provider/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

// GetRedisStore creates and returns a Redis storage instance
func GetRedisStore() *redis.Storage {
	return redis.New(redis.Config{
		Host: config.CacheConfig.CacheHost,
	})
}

// GetFromCache retrieves a value from Redis cache if caching is enabled
func GetFromCache(c *fiber.Ctx, key string) ([]byte, error) {
	// Check if caching is enabled
	if !config.CacheConfig.CacheEnabled {
		config.LogInfo(c, "Cache is disabled, skipping retrieval")
		return nil, nil
	}

	store := GetRedisStore()
	defer store.Close()

	val, err := store.Get(key)
	if err != nil {
		config.LogError(c, "Error retrieving from cache", "cache_key", key, "error", err.Error())
		return nil, err
	}

	if val != nil {
		config.LogInfo(c, "Cache hit", "cache_key", key)
	} else {
		config.LogInfo(c, "Cache miss", "cache_key", key)
	}

	return val, nil
}

// SetToCache stores a value in Redis cache with TTL if caching is enabled
func SetToCache(c *fiber.Ctx, key string, value []byte) error {
	// Check if caching is enabled
	if !config.CacheConfig.CacheEnabled {
		config.LogInfo(c, "Cache is disabled, skipping set")
		return nil
	}

	store := GetRedisStore()
	defer store.Close()

	// Convert TTL string to time.Duration (supports 5m, 1h, 30s, etc.)
	ttl := config.GetDurationFromEnv("CACHE_TTL", 5*time.Minute)

	if err := store.Set(key, value, ttl); err != nil {
		config.LogError(c, "Error setting cache", "cache_key", key, "ttl", config.CacheConfig.CacheTTL, "error", err.Error())
		return err
	}

	config.LogInfo(c, "Cache set", "cache_key", key, "ttl", config.CacheConfig.CacheTTL)
	return nil
}
