package repository

import (
	"errors"
)

var (
	// ErrFailedToReservePlayerSlot is the error returned when the player slot cannot be reserved
	ErrFailedToReservePlayerSlot = errors.New("failed to reserve player slot")

	// ErrFailedToUpdateTicket is the error returned when the ticket status cannot be updated
	ErrFailedToUpdateTicket = errors.New("failed to update ticket")

	// ErrFailedToGetTicket is the error returned when the ticket cannot be retrieved
	ErrFailedToGetTicket = errors.New("failed to get ticket")

	// ErrFailedToDeleteTicket is the error returned when the ticket cannot be deleted
	ErrFailedToDeleteTicket = errors.New("failed to delete ticket")

	// ErrTicketNotFound is the error returned when the match is not found
	ErrTicketNotFound = errors.New("ticket not found")

	// ErrFailedToEncodeTicket is the error returned when the ticket cannot be encoded
	ErrFailedToEncodeTicket = errors.New("failed to encode ticket")

	// ErrFailedToDecodeTicket is the error returned when the ticket cannot be decoded
	ErrFailedToDecodeTicket = errors.New("failed to decode ticket")

	// ErrSomethingOddHappened is the error returned when something odd happened
	ErrSomethingOddHappened = errors.New("something odd happened")
)
