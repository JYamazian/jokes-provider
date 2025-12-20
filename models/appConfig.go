package models

type AppConfig struct {
	// Server configuration
	Port        string
	Environment string
	AppHost     string

	// Logging configuration
	LogLevel         string
	LogFormat        string
	LogFormatType    string
	LogDisableColors string

	// Build information
	Version string
	Flavor  string

	// File paths
	JokesFilePath string

	// Request headers
	IPHeaderName      string
	CountryHeaderName string

	// Rate limiter configuration
	RateLimitEnabled     bool
	RateLimitMaxRequests int
	RateLimitDuration    string
}
