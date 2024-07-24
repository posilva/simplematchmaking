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
	srv := &MatchmakingService{
		logger:     logger,
		matchmaker: mm,
		repository: repo,
	}

	mm.Subscribe(srv)
	return srv
}

// HandleMatchResultError handles the match result error
func (s *MatchmakingService) HandleMatchResultError(err error) {
	s.logger.Error("Match result error received", err)
}

// HandleMatchResultOK handles the match result
func (s *MatchmakingService) HandleMatchResultOK(match domain.MatchResult) {
	now := time.Now().UTC().Unix()
	for _, t := range match.Tickets {
		s.logger.Info("Match result: Updating ticket", "ticketID", t.ID, "matchID", match.Match.ID)
		err := s.repository.UpdateTicket(context.Background(), domain.TicketStatus{
			ID:        t.ID,
			Timestamp: now,
			State:     domain.TicketStateMatched,
			MatchID:   match.Match.ID,
		})
		if err != nil {
			s.logger.Error("Failed to update ticket", err, "ticketID", t.ID, "matchID", match.Match.ID)
		}
	}
}

// FindMatch finds a match given a player
func (s *MatchmakingService) FindMatch(ctx context.Context, queue string, p domain.Player) (domain.Ticket, error) {
	ticketID := ksuid.New().String()
	now := time.Now().UTC().Unix()

	// check if the player is already in the queue
	ticketIDReserved, err := s.repository.ReservePlayerSlot(ctx, p.ID, queue, ticketID)
	if err != nil {
		s.logger.Error("Failed to reserve player slot", err)
		return domain.Ticket{}, fmt.Errorf("failed to reserve player slot: %v", err)
	}
	if ticketIDReserved != "" && ticketIDReserved != ticketID {
		s.logger.Info("Player already in the queue", "queue", queue, "player", p, "existingticketID", ticketIDReserved, "newticketID", ticketID)
		return domain.Ticket{
			ID: ticketIDReserved,
		}, nil
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
		Queue:     queue,
	}

	err = s.repository.UpdateTicket(ctx, status)
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
	ticket, err := s.repository.GetTicket(ctx, ticketID)
	if err != nil {
		s.logger.Error("Failed to get ticket status", err)
		return domain.Match{}, fmt.Errorf("failed to get ticket status: %v", err)
	}
	if ticket.State == domain.TicketStateMatched {
		return domain.Match{
			ID: ticket.MatchID,
		}, nil
	}
	return domain.Match{}, ErrMatchNotFound
}

// CancelMatch cancels a match given a ticket ID
func (s *MatchmakingService) CancelMatch(ctx context.Context, ticketID string) error {
	// if there is failure in the middle of the process, the slot will be stuck in the queue
	// for x amount of time as it should expire
	ticketStatus, err := s.repository.DeleteTicket(ctx, ticketID)
	if err != nil {
		s.logger.Error("Failed to delete ticket", err, "ticketID", ticketID)
		return fmt.Errorf("failed to delete ticket status: %v", err)
	}
	err = s.repository.DeletePlayerSlot(ctx, ticketStatus.PlayerID, ticketStatus.Queue)
	if err != nil {
		s.logger.Error("Failed to delete reservation", err, "ticketID", ticketID)
		return fmt.Errorf("failed to delete reservation: %v", err)
	}
	return nil
}
