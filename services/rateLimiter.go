package services

import (
	"jokes-provider/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRateLimiter() fiber.Handler {
	expiration := config.GetDurationFromEnv("RATE_LIMITER_EXPIRATION", 1*time.Minute)
	maxRequests := config.AppConfig.RateLimitMaxRequests

	config.LogInfo(nil, "Rate limiter initialized", "max_requests", maxRequests, "expiration", expiration)

	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			ip := c.Get(config.AppConfig.IPHeaderName)
			if ip == "" {
				ip = c.IP()
			}
			return ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			config.LogInfo(c, "Rate limit exceeded")
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	})
}
