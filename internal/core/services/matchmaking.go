package services

import (
	"context"
	"fmt"
	"time"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/segmentio/ksuid"
)

// MatchmakingService defines the Matchmaking service interface
type MatchmakingService struct {
	logger     ports.Logger
	matchmaker ports.Matchmaker
	repository ports.Repository
}

// NewMatchmakingService creates a new MatchmakingService
func NewMatchmakingService(logger ports.Logger, repo ports.Repository, mm ports.Matchmaker) *MatchmakingService {
	return &MatchmakingService{
		logger:     logger,
		matchmaker: mm,
		repository: repo,
	}
}

// FindMatch finds a match given a player
func (s *MatchmakingService) FindMatch(ctx context.Context, queue string, p domain.Player) (domain.Ticket, error) {
	ticketID := ksuid.New().String()
	now := time.Now().UTC().Unix()

	// check if the player is already in the queue
	ok, err := s.repository.ReservePlayerSlot(ctx, p.ID, queue, ticketID)
	if err != nil {
		s.logger.Error("Failed to reserve player slot", err)
		return domain.Ticket{}, fmt.Errorf("failed to reserve player slot: %v", err)
	}
	if !ok {
		s.logger.Info("Player already in the queue", "queue", queue, "player", p)
		return domain.Ticket{}, fmt.Errorf("player already in the queue")
	}

	err = s.matchmaker.AddPlayer(ctx, p)
	if err != nil {
		s.logger.Error("Failed to add player to the matchmaker", err)
		return domain.Ticket{}, fmt.Errorf("failed to add player to the matchmaker: %v", err)
	}

	status := domain.TicketStatus{
		ID:        ticketID,
		Timestamp: now,
		State:     domain.TicketStateQueued,
		PlayerID:  p.ID,
	}

	err = s.repository.UpdateTicketStatus(ctx, status)
	if err != nil {
		s.logger.Error("Failed to update ticket", err)
		return domain.Ticket{}, fmt.Errorf("failed to update ticket: %v", err)
	}

	return domain.Ticket{
		ID: ticketID,
	}, nil
}

// CheckMatch gets a match given a ticket ID
func (s *MatchmakingService) CheckMatch(ctx context.Context, ticketID string) (domain.Match, error) {
	_ = ticketID
	s.logger.Info("Match found", "ticketID", ticketID)
	return domain.Match{
		ID: "match1",
	}, nil
}

// CancelMatch cancels a match given a ticket ID
func (s *MatchmakingService) CancelMatch(ctx context.Context, ticketID string) error {
	_ = ticketID
	s.logger.Info("Match canceled", "ticketID", ticketID)
	return nil
}
