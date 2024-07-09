package services

import (
	"github.com/posilva/simplematchmaking/internal/core/domain"
)

// MatchmakingService defines the Matchmaking service interface
type MatchmakingService struct {
}

// NewMatchmakingService creates a new MatchmakingService
func NewMatchmakingService() *MatchmakingService {
	return &MatchmakingService{}
}

// FindMatch finds a match given a player
func (s *MatchmakingService) FindMatch(p domain.Player) (domain.Ticket, error) {
	_ = p
	return domain.Ticket{
		ID: "ticket1",
	}, nil
}
