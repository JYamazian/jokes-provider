package middleware

import (
	"crypto/tls"
	"crypto/x509"
	"jokes-provider/config"
	"jokes-provider/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
)

var redisStore *redis.Storage

func InitRedis() error {
	redisStore = redis.New(GetRedisConfig())
	return nil
}

func GetRedisStore() *redis.Storage {
	return redisStore
}

func CloseRedis() error {
	if redisStore != nil {
		return redisStore.Close()
	}
	return nil
}

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

func GetFromCache(c *fiber.Ctx, key string) ([]byte, error) {
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

func SetToCache(c *fiber.Ctx, key string, value []byte) error {
	store := GetRedisStore()

	// Convert TTL string to time.Duration (supports 5m, 1h, 30s, etc.)
	ttl := utils.GetDurationFromEnv(config.CacheConfig.CacheTTL, 5*time.Minute)

	if err := store.Set(key, value, ttl); err != nil {
		config.LogError(c, "Error setting cache", "cache_key", key, "ttl", config.CacheConfig.CacheTTL, "error", err.Error())
		return err
	}

	config.LogInfo(c, "Cache set", "cache_key", key, "ttl", config.CacheConfig.CacheTTL)
	return nil
}
