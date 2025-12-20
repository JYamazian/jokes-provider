package models

// Metadata represents the application metadata
type Metadata struct {
	App         AppInfo         `json:"app"`
	Server      ServerInfo      `json:"server"`
	Logging     LoggingInfo     `json:"logging"`
	Cache       CacheInfo       `json:"cache"`
	Files       FilesInfo       `json:"files"`
	Headers     HeadersInfo     `json:"headers"`
	RateLimiter RateLimiterInfo `json:"rate_limiter"`
	Fiber       FiberInfo       `json:"fiber"`
}

type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

type ServerInfo struct {
	Port        string `json:"port"`
	Environment string `json:"environment"`
	Timestamp   string `json:"timestamp"`
}

type LoggingInfo struct {
	Level         string `json:"level"`
	Format        string `json:"format"`
	FormatType    string `json:"format_type"`
	DisableColors string `json:"disable_colors"`
}

type CacheInfo struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	TTL     string `json:"ttl"`
}

type FilesInfo struct {
	JokesPath string `json:"jokes_path"`
}

type HeadersInfo struct {
	IPHeaderName      string `json:"ip_header_name"`
	CountryHeaderName string `json:"country_header_name"`
}

type RateLimiterInfo struct {
	MaxRequests int    `json:"max_requests"`
	Duration    string `json:"duration"`
}

type FiberInfo struct {
	Prefork       bool `json:"prefork"`
	CaseSensitive bool `json:"case_sensitive"`
	StrictRouting bool `json:"strict_routing"`
}
