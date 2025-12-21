package utils

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 100 // Default to 100 requests
	}
	return val
}

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
