package services

import (
	"context"
	"time"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
)

// Matchmaker is the Matchmaker implementation using
type Matchmaker struct {
	queue              ports.Queue
	config             domain.MatchmakerConfig
	scheduler          *Scheduler
	logger             ports.Logger
	matchResultHandler ports.MatchResultsListHandler
	allBrackets        []string
}

// NewMatchmaker creates a new Matchmaker
func NewMatchmaker(queue ports.Queue, cfg domain.MatchmakerConfig, logger ports.Logger) (*Matchmaker, error) {
	mm := &Matchmaker{
		queue:       queue,
		config:      cfg,
		logger:      logger,
		allBrackets: make([]string, 0),
	}

	// Create a new scheduler
	mm.scheduler = NewScheduler(cfg.IntervalSecs, mm.Matchmake)
	return mm, nil
}

// AddPlayer adds a player to the matchmaker
func (m *Matchmaker) AddPlayer(ctx context.Context, ticketID string, p domain.Player) error {
	qe := domain.QueueEntry{
		TicketID: ticketID,
		PlayerID: p.ID,
		Ranking:  p.Ranking,
	}

	return m.queue.Enqueue(ctx, qe)
}

// Matchmake finds a match
func (m *Matchmaker) Matchmake() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.MakeTimeoutSecs)*time.Second)
	defer func() {
		cancel()
	}()

	mr, err := m.queue.Make(ctx)
	if err != nil {
		m.logger.Error("failed to make match: %v", err, m.queue.Name())
		m.matchResultHandler.HandleMatchResultsError(m.queue.Name(), err)
		return
	}
	if m.matchResultHandler != nil {
		m.matchResultHandler.HandleMatchResultsOK(m.queue.Name(), mr)
	}
}

// Subscribe subscribes to match results
func (m *Matchmaker) Subscribe(handler ports.MatchResultsListHandler) {
	m.matchResultHandler = handler
}
