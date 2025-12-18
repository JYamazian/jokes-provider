package config

import (
	"jokes-provider/models"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var AppConfig *models.AppConfig
var CacheConfig *models.CacheConfig

// LoadEnvVars loads environment variables from OS and initializes the global Config
func LoadEnvVars() {
	AppConfig = &models.AppConfig{
		// Server configuration with defaults
		Port:        getEnv("PORT", "3000"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Logging
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		LogFormat:        getEnv("LOG_FORMAT", "[${ip}]:${port} ${status} - ${method} ${path}"),
		LogFormatType:    getEnv("LOG_FORMAT_TYPE", "text"),
		LogDisableColors: getEnv("LOG_DISABLE_COLORS", "false"),

		// Build information (loaded from environment, set by Docker build args)
		Version: getEnv("BUILD_VERSION", "dev"),
		Flavor:  getEnv("BUILD_FLAVOR", "development"),

		// Fiber configuration
		FiberConfig: models.FiberConfig{
			Prefork:       getEnv("FIBER_PREFORK", "false") == "true",
			CaseSensitive: getEnv("FIBER_CASE_SENSITIVE", "false") == "true",
			StrictRouting: getEnv("FIBER_STRICT_ROUTING", "false") == "true",
		},

		// File paths
		JokesFilePath: getEnv("JOKES_FILE_PATH", "/data/jokes.csv"),

		// Request headers
		IPHeaderName:      getEnv("IP_HEADER_NAME", "X-Forwarded-For"),
		CountryHeaderName: getEnv("COUNTRY_HEADER_NAME", "X-Country-Name"),

		// Rate limiter configuration
		RateLimitMaxRequests: parseInt(getEnv("RATE_LIMIT_MAX_REQUESTS", "100")),
		RateLimitDuration:    getEnv("RATE_LIMITER_EXPIRATION", "1m"),
	}

	CacheConfig = &models.CacheConfig{
		CacheURL:     getEnv("CACHE_URL", "localhost"),
		CacheEnabled: getEnv("CACHE_ENABLED", "true") == "true",
		CacheTTL:     getEnv("CACHE_TTL", "5m"),
	}
}

// InitializeApp creates and initializes the Fiber app
func InitializeApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       AppConfig.FiberConfig.Prefork,
		CaseSensitive: AppConfig.FiberConfig.CaseSensitive,
		StrictRouting: AppConfig.FiberConfig.StrictRouting,
		ServerHeader:  "Go Fiber - Jokes Provider",
		AppName:       "Jokes Provider API",
	})

	// Initialize logger
	InitializeLogger(app)

	return app
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt converts a string to int with a default value
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 100 // Default to 100 requests
	}
	return val
}

// GetDurationFromEnv retrieves a duration from environment variable with a fallback default
func GetDurationFromEnv(envVar string, defaultDuration time.Duration) time.Duration {
	durationStr := os.Getenv(envVar)
	if durationStr == "" {
		return defaultDuration
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return defaultDuration
	}
	return duration
}
