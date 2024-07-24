package services

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/segmentio/ksuid"
)

// Matchmaker is the Matchmaker implementation using
type Matchmaker struct {
	queue              ports.Queue
	config             domain.MatchmakerConfig
	scheduler          *Scheduler
	logger             ports.Logger
	matchResultHandler ports.MatchResultHandler
}

// NewMatchmaker creates a new Matchmaker
func NewMatchmaker(queue ports.Queue, cfg domain.MatchmakerConfig, logger ports.Logger) (*Matchmaker, error) {
	mm := &Matchmaker{
		queue:  queue,
		config: cfg,
		logger: logger,
	}
	mm.scheduler = NewScheduler(cfg.IntervalSecs, mm.Matchmake)
	return mm, nil
}

// AddPlayer adds a player to the matchmaker
func (m *Matchmaker) AddPlayer(ctx context.Context, p domain.Player) error {
	return m.queue.AddPlayer(ctx, p)
}

// Matchmake finds a match
func (m *Matchmaker) Matchmake() {
	ctx := context.Background()
	matchID := ksuid.New().String()
	mr, err := m.queue.Make(ctx, matchID)
	if err != nil {
		m.logger.Error("failed to make match: %v", err, m.queue.Name())
		m.matchResultHandler.HandleMatchResultError(err)
		return
	}
	if m.matchResultHandler != nil {
		m.matchResultHandler.HandleMatchResultOK(mr)
	}
}

// Subscribe subscribes to match results
func (m *Matchmaker) Subscribe(handler ports.MatchResultHandler) {
	m.matchResultHandler = handler
}
