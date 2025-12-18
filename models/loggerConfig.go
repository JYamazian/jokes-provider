package models

import "log"

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string `json:"timestamp,omitempty"`
	Level     string `json:"level"`
	RequestID string `json:"request_id,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	Country   string `json:"country,omitempty"`
	Message   string `json:"message"`
}

// ContextLogger provides context-based logging
type ContextLogger struct {
	Logger *log.Logger
	Format string
}
