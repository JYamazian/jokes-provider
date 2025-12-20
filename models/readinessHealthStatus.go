package models

// HealthStatus represents the health check result
type ReadinessHealthStatus struct {
	Ready  bool   `json:"ready"`
	Reason string `json:"reason,omitempty"`
	Redis  string `json:"redis,omitempty"`
	CSV    string `json:"csv,omitempty"`
}
