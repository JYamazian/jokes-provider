package config

import (
	"jokes-provider/models"
	"jokes-provider/utils"
)

var AppConfig *models.AppConfig
var CacheConfig *models.CacheConfig
var FiberConfig *models.FiberConfig

// LoadEnvVars loads environment variables from OS and initializes the global Config
func LoadEnvVars() {
	AppConfig = &models.AppConfig{
		// Server configuration with defaults
		Port:        utils.GetEnv("PORT", "3000"),
		Environment: utils.GetEnv("ENVIRONMENT", "development"),

		// Logging
		LogLevel:         utils.GetEnv("LOG_LEVEL", "info"),
		LogFormat:        utils.GetEnv("LOG_FORMAT", "[${ip}]:${port} ${status} - ${method} ${path}"),
		LogFormatType:    utils.GetEnv("LOG_FORMAT_TYPE", "text"),
		LogDisableColors: utils.GetEnv("LOG_DISABLE_COLORS", "false"),

		// Build information (loaded from environment, set by Docker build args)
		Version: utils.GetEnv("BUILD_VERSION", "dev"),
		Flavor:  utils.GetEnv("BUILD_FLAVOR", "development"),
		// File paths
		JokesFilePath: utils.GetEnv("JOKES_FILE_PATH", "/data/jokes.csv"),
		// Request headers
		IPHeaderName:      utils.GetEnv("IP_HEADER_NAME", "X-Forwarded-For"),
		CountryHeaderName: utils.GetEnv("COUNTRY_HEADER_NAME", "X-Country-Name"),

		// Rate limiter configuration
		RateLimitEnabled:     utils.GetEnv("RATE_LIMIT_ENABLED", "false") == "true",
		RateLimitMaxRequests: utils.ParseInt(utils.GetEnv("RATE_LIMIT_MAX_REQUESTS", "100")),
		RateLimitDuration:    utils.GetEnv("RATE_LIMITER_EXPIRATION", "1m"),
	}

	// Cache configuration
	CacheConfig = &models.CacheConfig{
		CacheURL:            utils.GetEnv("CACHE_URL", "localhost"),
		CacheEnabled:        utils.GetEnv("CACHE_ENABLED", "true") == "true",
		CacheTTL:            utils.GetEnv("CACHE_TTL", "5m"),
		CacheCaCertPath:     utils.GetEnv("CACHE_CA_CERT", ""),
		CacheClientCertPath: utils.GetEnv("CACHE_CLIENT_CERT", ""),
		CacheClientKeyPath:  utils.GetEnv("CACHE_CLIENT_KEY", ""),
	}

	// Fiber configuration
	FiberConfig = &models.FiberConfig{
		Prefork:       utils.GetEnv("FIBER_PREFORK", "false") == "true",
		CaseSensitive: utils.GetEnv("FIBER_CASE_SENSITIVE", "false") == "true",
		StrictRouting: utils.GetEnv("FIBER_STRICT_ROUTING", "false") == "true",
	}
}
