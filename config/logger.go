package config

import (
	"encoding/json"
	"fmt"
	"jokes-provider/models"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// ContextLogger wraps models.ContextLogger to allow defining methods
type ContextLogger struct {
	*models.ContextLogger
}

// contextLogger is the global logger instance
var contextLogger *ContextLogger

// NewContextLogger creates a new context logger
func NewContextLogger(format string) *ContextLogger {
	return &ContextLogger{
		ContextLogger: &models.ContextLogger{
			Logger: log.New(os.Stdout, "", 0),
			Format: format,
		},
	}
}

// LogWithContext logs a message with Fiber context information
func (cl *ContextLogger) LogWithContext(c *fiber.Ctx, level string, message string, fields ...interface{}) {
	entry := &models.LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	// Only access context if it's not nil
	if c != nil {
		// Check for forwarded IP first, then fall back to direct IP
		if forwardedFor := c.Get(AppConfig.IPHeaderName); forwardedFor != "" {
			entry.IPAddress = forwardedFor
		} else {
			entry.IPAddress = c.IP()
		}

		// Check for country name header
		if country := c.Get(AppConfig.CountryHeaderName); country != "" {
			entry.Country = country
		}

		if requestID := c.Get(fiber.HeaderXRequestID); requestID != "" {
			entry.RequestID = requestID
		}
	}

	if cl.Format == "json" {
		cl.logJSON(entry, fields...)
	} else {
		cl.printTextFormat(entry, fields...)
	}
}

// logJSON logs entry as JSON with additional structured fields
func (cl *ContextLogger) logJSON(entry *models.LogEntry, fields ...interface{}) {
	logMap := map[string]interface{}{
		"timestamp": entry.Timestamp,
		"level":     entry.Level,
		"message":   entry.Message,
		"version":   fmt.Sprintf("%s-%s", AppConfig.Flavor, AppConfig.Version),
		"pid":       os.Getpid(),
	}

	if entry.RequestID != "" {
		logMap["request_id"] = entry.RequestID
	}
	if entry.IPAddress != "" {
		logMap["ip_address"] = entry.IPAddress
	}
	if entry.Country != "" {
		logMap["country"] = entry.Country
	}

	// Add custom fields (key, value, key, value...)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprintf("%v", fields[i])
			value := fields[i+1]
			logMap[key] = value
		}
	}

	jsonData, _ := json.Marshal(logMap)
	cl.Logger.Println(string(jsonData))
}

// printTextFormat formats and prints log entry as text with structured fields
func (cl *ContextLogger) printTextFormat(entry *models.LogEntry, fields ...interface{}) {
	logMsg := fmt.Sprintf("[%s] [%s]", entry.Timestamp, entry.Level)
	if entry.RequestID != "" {
		logMsg += fmt.Sprintf(" [%s]", entry.RequestID)
	}
	if entry.IPAddress != "" {
		logMsg += fmt.Sprintf(" [%s]", entry.IPAddress)
	}
	if entry.Country != "" {
		logMsg += fmt.Sprintf(" [%s]", entry.Country)
	}
	logMsg += fmt.Sprintf(" %s", entry.Message)

	// Add build info and structured fields to text format
	if AppConfig.Version != "" || AppConfig.Flavor != "" || len(fields) > 0 {
		logMsg += " |"
		if AppConfig.Version != "" && AppConfig.Flavor != "" {
			logMsg += fmt.Sprintf(" version=%s-%s", AppConfig.Flavor, AppConfig.Version)
		}
		for i := 0; i < len(fields); i += 2 {
			if i+1 < len(fields) {
				key := fmt.Sprintf("%v", fields[i])
				value := fmt.Sprintf("%v", fields[i+1])
				logMsg += fmt.Sprintf(" %s=%s", key, value)
			}
		}
	}

	cl.Logger.Println(logMsg)
}

// InitializeLogger configures logging middleware for the Fiber app
func InitializeLogger(app *fiber.App) {
	app.Use(requestid.New())

	logFormat := getEnv("LOG_FORMAT_TYPE", "text")
	contextLogger = NewContextLogger(logFormat)

	// Setup Fiber healthcheck middleware for liveness probe
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/health/liveness",
	}))

	// Setup Fiber logger middleware
	if logFormat == "json" {
		app.Use(logger.New(logger.Config{
			Format: "${time_rfc3339} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
			Output: os.Stdout,
		}))
	} else {
		app.Use(logger.New(logger.Config{
			Format: "[${time_rfc3339}] ${status} ${latency} ${ip} ${method} ${path}${error}\n",
			Output: os.Stdout,
		}))
	}
}

// LogInfo logs an info message with Fiber context
func LogInfo(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "INFO", message, fields...)
}

// LogError logs an error message with Fiber context
func LogError(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "ERROR", message, fields...)
}

// LogDebug logs a debug message with Fiber context
func LogDebug(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "DEBUG", message, fields...)
}

// LogStartupInfo logs startup information using Fiber's logger
func LogStartupInfo(version, flavor string) {
	// Use the context logger to log startup info with structured fields
	LogInfo(nil, "Application started", "version", version, "flavor", flavor, "environment", AppConfig.Environment, "port", AppConfig.Port)
}
