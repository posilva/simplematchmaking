package ports

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
)

// Codec defines the interface to encode/decode data
type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// MatchResultsListHandler defines the interface to handle match results
type MatchResultsListHandler interface {
	HandleMatchResultsOK(matches []domain.MatchResult)
	HandleMatchResultsError(err error)
}

// Matchmaker defines the Matchmaker interface
type Matchmaker interface {
	Matchmake()
	AddPlayer(ctx context.Context, p domain.Player) error
	Subscribe(handler MatchResultsListHandler)
}

// Queue defines the Queue interface
type Queue interface {
	Enqueue(ctx context.Context, qe domain.QueueEntry) error
	Make(ctx context.Context) ([]domain.MatchResult, error)
	Name() string
}

// MatchmakingService defines the matchmaking service interface
type MatchmakingService interface {
	FindMatch(ctx context.Context, queue string, p domain.Player) (domain.Ticket, error)
	CheckMatch(ctx context.Context, ticketID string) (domain.Match, error)
	CancelMatch(ctx context.Context, ticketID string) error
}

// Repository defines the interface to handle with data to be stored/retrieved
type Repository interface {
	// ReservePlayerSlot reserves a player slot in the queue
	ReservePlayerSlot(ctx context.Context, playerID string, slot string, ticketID string) (string, error)
	// DeletePlayerSlot deletes a player slot in the queue
	DeletePlayerSlot(ctx context.Context, playerID string, slot string) error
	// UpdateTicket updates the ticket status
	UpdateTicket(ctx context.Context, ticket domain.TicketRecord) error
	// GetTicket gets the ticket status
	GetTicket(ctx context.Context, ticketID string) (domain.TicketRecord, error)
	// DeleteTicket deletes the ticket status
	DeleteTicket(ctx context.Context, ticketID string) (domain.TicketRecord, error)
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
