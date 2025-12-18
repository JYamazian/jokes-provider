package helpers

import (
	"jokes-provider/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

// CheckRedisStatus checks if Redis is accessible
func CheckRedisStatus(c *fiber.Ctx) bool {
	store := redis.New(redis.Config{
		Host: config.CacheConfig.CacheHost,
	})
	defer store.Close()

	// Try to set and get a test key
	testKey := "health_check"
	if err := store.Set(testKey, []byte("ok"), 5*time.Second); err != nil {
		config.LogError(c, "Redis health check failed on SET", "error", err.Error())
		return false
	}

	val, err := store.Get(testKey)
	if err != nil {
		config.LogError(c, "Redis health check failed on GET", "error", err.Error())
		return false
	}

	if string(val) != "ok" {
		config.LogError(c, "Redis health check failed: invalid response")
		return false
	}

	config.LogInfo(c, "Redis health check passed")
	return true
}
