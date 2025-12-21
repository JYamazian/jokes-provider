package services

import (
	"jokes-provider/config"
	"jokes-provider/models"
	"time"
)

// MetadataService handles metadata business logic
type MetadataService struct{}

// NewMetadataService creates a new MetadataService instance
func NewMetadataService() *MetadataService {
	return &MetadataService{}
}

// GetMetadata returns the application metadata
func (s *MetadataService) GetMetadata() models.Metadata {
	return models.Metadata{
		App: models.AppInfo{
			Name:    "Jokes Provider API",
			Version: config.AppConfig.Version,
			Flavor:  config.AppConfig.Flavor,
		},
		Server: models.ServerInfo{
			Port:        config.AppConfig.Port,
			Environment: config.AppConfig.Environment,
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		Logging: models.LoggingInfo{
			Level:         config.AppConfig.LogLevel,
			Format:        config.AppConfig.LogFormat,
			FormatType:    config.AppConfig.LogFormatType,
			DisableColors: config.AppConfig.LogDisableColors,
		},
		Cache: models.CacheInfo{
			Enabled: config.CacheConfig.CacheEnabled,
			URL:     config.CacheConfig.CacheURL,
			TTL:     config.CacheConfig.CacheTTL,
		},
		Files: models.FilesInfo{
			JokesPath: config.AppConfig.JokesFilePath,
		},
		Headers: models.HeadersInfo{
			IPHeaderName:      config.AppConfig.IPHeaderName,
			CountryHeaderName: config.AppConfig.CountryHeaderName,
		},
		RateLimiter: models.RateLimiterInfo{
			Enabled:     config.AppConfig.RateLimitEnabled,
			MaxRequests: config.AppConfig.RateLimitMaxRequests,
			Duration:    config.AppConfig.RateLimitDuration,
		},
		Fiber: models.FiberInfo{
			Prefork:       config.FiberConfig.Prefork,
			CaseSensitive: config.FiberConfig.CaseSensitive,
			StrictRouting: config.FiberConfig.StrictRouting,
		},
	}
}
