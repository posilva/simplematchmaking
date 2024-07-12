package services

import (
	"github.com/posilva/simplematchmaking/internal/adapters/output/logging"
	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
)

// MatchmakingService defines the Matchmaking service interface
type MatchmakingService struct {
	logger       ports.Logger
	queueManager ports.QueueManager
	repository   ports.Repository
}

// NewMatchmakingService creates a new MatchmakingService
func NewMatchmakingService() *MatchmakingService {
	return &MatchmakingService{
		logger: logging.NewSimpleLogger(),
	}
}

// FindMatch finds a match given a player
func (s *MatchmakingService) FindMatch(queue string, p domain.Player) (domain.Ticket, error) {
	_ = p
	_ = queue
	s.logger.Info("Match found", "queue", queue, "player", p)
	return domain.Ticket{
		ID: "ticket1",
	}, nil
}

// GetMatch gets a match given a ticket ID
func (s *MatchmakingService) GetMatch(ticketID string) (domain.Match, error) {
	_ = ticketID
	s.logger.Info("Match found", "ticketID", ticketID)
	return domain.Match{
		ID: "match1",
	}, nil
}

// CancelMatch cancels a match given a ticket ID
func (s *MatchmakingService) CancelMatch(ticketID string) error {
	_ = ticketID
	s.logger.Info("Match canceled", "ticketID", ticketID)
	return nil
}
