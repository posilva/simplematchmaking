package services

import "fmt"

var (
	// ErrFailedToReservePlayerSlot is the error returned when the player slot cannot be reserved
	ErrFailedToReservePlayerSlot = fmt.Errorf("failed to reserve player slot")

	// ErrPlayerAlreadyInQueue is the error returned when the player is already in the queue
	ErrPlayerAlreadyInQueue = fmt.Errorf("player already in the queue")

	// ErrFailedToAddPlayerToMatchmaker is the error returned when the player cannot be added to the matchmaker
	ErrFailedToAddPlayerToMatchmaker = fmt.Errorf("failed to add player to the matchmaker")

	// ErrFailedToUpdateTicket is the error returned when the ticket status cannot be updated
	ErrFailedToUpdateTicket = fmt.Errorf("failed to update ticket")

	// ErrFailedToGetTicketStatus is the error returned when the ticket status cannot be retrieved
	ErrFailedToGetTicketStatus = fmt.Errorf("failed to get ticket status")

	// ErrMatchNotFound is the error returned when the match is not found
	ErrMatchNotFound = fmt.Errorf("match not found")
)
