package ports

import "github.com/posilva/simplematchmaking/internal/core/domain"

// Repository defines the interface to handle with
type Repository interface {
}

// Logger defines a basic logger interface
type Logger interface {
	Debug(msg string, v ...interface{}) error
	Info(msg string, v ...interface{}) error
	Error(msg string, v ...interface{}) error
}

// MatchmakingService defines the matchmaking service interface
type MatchmakingService interface {
	AddPlayer(domain.Player) error
}

// Provider generic interface
type Provider[T any] interface {
	Provide() (T, error)
}

// TelemetryReporter defines the interface to report metrics
type TelemetryReporter interface {
	SetDefaultTags(tags map[string]string)
	ReportGauge(name string, value float64, tags map[string]string)
	ReportCounter(name string, value float64, tags map[string]string)
	ReportHistogram(name string, value float64, tags map[string]string)
	ReportSummary(name string, value float64, tags map[string]string)
}
