package services

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
)

// Matchmaker is the Matchmaker implementation using
type Matchmaker struct {
	queue  ports.Queue
	config domain.MatchmakerConfig
}

// NewMatchmaker creates a new Matchmaker
func NewMatchmaker(queue ports.Queue, cfg domain.MatchmakerConfig) (*Matchmaker, error) {
	return &Matchmaker{
		queue:  queue,
		config: cfg,
	}, nil
}

// AddPlayer adds a player to the matchmaker
func (m *Matchmaker) AddPlayer(ctx context.Context, p domain.Player) error {
	return m.queue.AddPlayer(ctx, p)
}

// Match finds a match
func (m *Matchmaker) Match(ctx context.Context) (domain.MatchResult, error) {
	return domain.MatchResult{
		Match: domain.Match{
			ID: "match1",
		},
		Tickets: []domain.Ticket{{ID: "ticket1"}, {ID: "ticket2"}},
	}, nil
}
