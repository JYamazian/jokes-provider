package config

import (
	"encoding/json"
	"fmt"
	"jokes-provider/models"
	"jokes-provider/utils"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type ContextLogger struct {
	*models.ContextLogger
}

var contextLogger *ContextLogger

func NewContextLogger(format string) *ContextLogger {
	return &ContextLogger{
		ContextLogger: &models.ContextLogger{
			Logger: log.New(os.Stdout, "", 0),
			Format: format,
		},
	}
}

func (cl *ContextLogger) LogWithContext(c *fiber.Ctx, level string, message string, fields ...interface{}) {
	entry := &models.LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	if c != nil {
		if forwardedFor := c.Get(AppConfig.IPHeaderName); forwardedFor != "" {
			entry.IPAddress = forwardedFor
		} else {
			entry.IPAddress = c.IP()
		}

		if country := c.Get(AppConfig.CountryHeaderName); country != "" {
			entry.Country = country
		}

		if requestID := c.Get(fiber.HeaderXRequestID); requestID != "" {
			entry.RequestID = requestID
		}
	}

	if cl.ContextLogger.Format == "json" {
		cl.logJSON(entry, fields...)
	} else {
		cl.printTextFormat(entry, fields...)
	}
}

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
	cl.ContextLogger.Logger.Println(string(jsonData))
}

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

	cl.ContextLogger.Logger.Println(logMsg)
}

func InitializeLogger(app *fiber.App) {
	app.Use(requestid.New())

	logFormat := utils.GetEnv("LOG_FORMAT_TYPE", "text")
	contextLogger = NewContextLogger(logFormat)

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

func LogInfo(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "INFO", message, fields...)
}

func LogError(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "ERROR", message, fields...)
}

func LogDebug(c *fiber.Ctx, message string, fields ...interface{}) {
	contextLogger.LogWithContext(c, "DEBUG", message, fields...)
}

func LogStartupInfo(version, flavor string) {
	LogInfo(nil, "Application started", "version", version, "flavor", flavor, "environment", AppConfig.Environment, "port", AppConfig.Port)
}
