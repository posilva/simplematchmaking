package ports

import "github.com/posilva/simplematchmaking/internal/core/domain"

// Queue defines the Matchmaking queue interface
type Queue interface {
	Enqueue(domain.Ticket) error
	Dequeue(num int) (domain.Ticket, error)
}

// QueueManager defines the Matchmaking queue manager interface
type QueueManager interface {
	Enqueue(name string, t domain.Ticket) error
	Dequeue(name string, num int) (domain.Ticket, error)
}

// MatchmakingService defines the matchmaking service interface
type MatchmakingService interface {
	FindMatch(queue string, p domain.Player) (domain.Ticket, error)
	GetMatch(ticketID string) (domain.Match, error)
	CancelMatch(ticketID string) error
}

// Repository defines the interface to handle with
type Repository interface {
}

// Logger defines a basic logger interface
type Logger interface {
	Debug(msg string, v ...interface{}) error
	Info(msg string, v ...interface{}) error
	Error(msg string, v ...interface{}) error
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
