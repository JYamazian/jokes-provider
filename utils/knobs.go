package utils

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// getEnv retrieves an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt converts a string to int with a default value
func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 100 // Default to 100 requests
	}
	return val
}

// GetDurationFromEnv retrieves a duration from config environment variable with a fallback default
func GetDurationFromEnv(envVar string, defaultDuration time.Duration) time.Duration {
	if envVar == "" {
		return defaultDuration
	}

	duration, err := time.ParseDuration(envVar)
	if err != nil {
		return defaultDuration
	}
	return duration
}

// ShouldSkipCache checks if the Cache-Control header contains "no-cache"
func ShouldSkipCache(c *fiber.Ctx) bool {
	cacheControl := c.Get(HeaderCacheControl)
	return strings.Contains(strings.ToLower(cacheControl), CacheControlNoCache)
}
