package services

import "github.com/posilva/simplematchmaking/internal/core/domain"

// MatchmakingService defines the Matchmaking service interface
type MatchmakingService struct {
}

// NewMatchmakingService creates a new MatchmakingService
func NewMatchmakingService() *MatchmakingService {
	return &MatchmakingService{}
}

// AddPlayer adds a new player
func (s *MatchmakingService) AddPlayer(p domain.Player) error {
	_ = p
	return nil
}
