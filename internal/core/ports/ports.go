package ports

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
)

// Matchmaker defines the Matchmaker interface
type Matchmaker interface {
	Match(ctx context.Context) (domain.MatchResult, error)
	AddPlayer(ctx context.Context, p domain.Player) error
}

// MatchmakingService defines the matchmaking service interface
type MatchmakingService interface {
	FindMatch(ctx context.Context, queue string, p domain.Player) (domain.Ticket, error)
	CheckMatch(ctx context.Context, ticketID string) (domain.Match, error)
	CancelMatch(ctx context.Context, ticketID string) error
}

// Repository defines the interface to handle with
type Repository interface {
	ReservePlayerSlot(ctx context.Context, playerID string, slot string, ticketID string) (bool, error)
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
