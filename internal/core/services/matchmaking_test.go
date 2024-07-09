package services

import (
	"testing"

	"github.com/posilva/simplematchmaking/internal/core/domain"
)

func TestMatchmakingService_FindMatch(t *testing.T) {
	s := NewMatchmakingService()
	ticket, err := s.FindMatch(domain.Player{
		ID: "player1",
	})
	if err != nil {
		t.Errorf("failed to find a match: %v", err)
	}
	if ticket.ID == "" {
		t.Errorf("failed to find a match")
	}
}
