package middleware

import (
	"crypto/tls"
	"crypto/x509"
	"jokes-provider/config"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

var redisStore *redis.Storage

// GetRedisConfig creates and returns a Redis storage configuration with TLS support
func GetRedisConfig() redis.Config {
	cfg := redis.Config{
		URL: config.CacheConfig.CacheURL,
	}

	if config.CacheConfig.CacheCaCertPath != "" || config.CacheConfig.CacheClientCertPath != "" {
		tlsConfig := &tls.Config{}

		// Load CA certificate
		if config.CacheConfig.CacheCaCertPath != "" {
			caCert, err := os.ReadFile(config.CacheConfig.CacheCaCertPath)
			if err != nil {
				config.LogError(nil, "Failed to read Redis CA cert", "error", err.Error())
			} else {
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				tlsConfig.RootCAs = caCertPool
			}
		}

		// Load client certificate and key
		if config.CacheConfig.CacheClientCertPath != "" && config.CacheConfig.CacheClientKeyPath != "" {
			clientCert, err := tls.LoadX509KeyPair(config.CacheConfig.CacheClientCertPath, config.CacheConfig.CacheClientKeyPath)
			if err != nil {
				config.LogError(nil, "Failed to load Redis client cert/key", "error", err.Error())
			} else {
				tlsConfig.Certificates = []tls.Certificate{clientCert}
			}
		}

		cfg.TLSConfig = tlsConfig
	}

	return cfg
}

// InitRedis initializes the Redis connection once at app startup
func InitRedis() error {
	redisStore = redis.New(GetRedisConfig())
	return nil
}

// GetRedisStore returns the already-initialized Redis connection (singleton)
func GetRedisStore() *redis.Storage {
	return redisStore
}

// CloseRedis closes the Redis connection (call on app shutdown)
func CloseRedis() error {
	if redisStore != nil {
		return redisStore.Close()
	}
	return nil
}

// GetFromCache retrieves a value from Redis cache if caching is enabled
func GetFromCache(c *fiber.Ctx, key string) ([]byte, error) {
	// Check if caching is enabled
	if !config.CacheConfig.CacheEnabled {
		config.LogInfo(c, "Cache is disabled, skipping retrieval")
		return nil, nil
	}

	store := GetRedisStore()

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

	// Convert TTL string to time.Duration (supports 5m, 1h, 30s, etc.)
	ttl := config.GetDurationFromEnv("CACHE_TTL", 5*time.Minute)

	if err := store.Set(key, value, ttl); err != nil {
		config.LogError(c, "Error setting cache", "cache_key", key, "ttl", config.CacheConfig.CacheTTL, "error", err.Error())
		return err
	}

	config.LogInfo(c, "Cache set", "cache_key", key, "ttl", config.CacheConfig.CacheTTL)
	return nil
}
