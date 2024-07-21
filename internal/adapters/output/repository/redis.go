// Package repository provides the repository implementation for the output port
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/posilva/simplematchmaking/internal/core/ports"
	"github.com/redis/rueidis"
)

var (
	// TODO: move to config
	reservationTimeEx = int64(60)
	redisCallTimeout  = 1000 * time.Millisecond
)

// RedisRepository is the Redis Repository
type RedisRepository struct {
	client rueidis.Client
	codec  ports.Codec
}

// NewRedisRepository creates a new RedisRepository
func NewRedisRepository(client rueidis.Client, codec ports.Codec) *RedisRepository {
	return &RedisRepository{
		client: client,
		codec:  codec,
	}
}

// ReservePlayerSlot reserves a player slot in the queue
func (r *RedisRepository) ReservePlayerSlot(ctx context.Context, playerID string, slot string, ticketID string) (bool, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.playerSlotKey(slot, playerID)
	value := "status:reserved:ticket:" + ticketID
	cmd := r.client.B().Set().Key(key).Value(value).Nx().ExSeconds(reservationTimeEx).Build()
	resp, err := r.client.Do(ctxWithTimeout, cmd).AsBool()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return resp, fmt.Errorf("failed to reserve player slot: '%v'  %v", slot, err)
		}
	}
	return resp, nil
}

// DeletePlayerSlot deletes a player slot in the queue
func (r *RedisRepository) DeletePlayerSlot(ctx context.Context, playerID string, slot string) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.playerSlotKey(slot, playerID)
	cmd := r.client.B().Del().Key(key).Build()
	return r.client.Do(ctxWithTimeout, cmd).Error()
}

// UpdateTicket updates the ticket status
func (r *RedisRepository) UpdateTicket(ctx context.Context, status domain.TicketStatus) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.ticketKey(status.ID)
	value, err := r.codec.Encode(status)
	if err != nil {
		return fmt.Errorf("failed to encode ticket %v", err)
	}

	// TODO: this should have a TTL
	cmd := r.client.B().Set().Key(key).Value(string(value)).Build()
	err = r.client.Do(ctxWithTimeout, cmd).Error()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return fmt.Errorf("failed to update ticket %v", err)
		}
		return fmt.Errorf("ticket not found: %v", err)
	}
	return nil
}

// GetTicket gets the ticket status
func (r *RedisRepository) GetTicket(ctx context.Context, ticketID string) (domain.TicketStatus, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.ticketKey(ticketID)
	cmd := r.client.B().Get().Key(key).Build()
	resp, err := r.client.Do(ctxWithTimeout, cmd).AsBytes()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return domain.TicketStatus{}, fmt.Errorf("failed to get ticket %v", err)
		}
		return domain.TicketStatus{}, fmt.Errorf("ticket not found: %v", err)
	}

	var status domain.TicketStatus
	err = r.codec.Decode([]byte(resp), &status)
	if err != nil {
		return domain.TicketStatus{}, fmt.Errorf("failed to decode ticket %v", err)
	}
	return status, nil
}

// DeleteTicket deletes the ticket status
func (r *RedisRepository) DeleteTicket(ctx context.Context, ticketID string) (domain.TicketStatus, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.ticketKey(ticketID)
	cmd := r.client.B().Getdel().Key(key).Build()
	b, err := r.client.Do(ctxWithTimeout, cmd).AsBytes()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return domain.TicketStatus{}, fmt.Errorf("failed to delete ticket %v", err)
		}
		return domain.TicketStatus{}, fmt.Errorf("ticket not found: %v", err)
	}
	var status domain.TicketStatus
	err = r.codec.Decode(b, &status)
	if err != nil {
		return domain.TicketStatus{}, fmt.Errorf("failed to decode ticket %v", err)
	}
	return status, nil
}

// playerSlotKey returns the key for the player slot to use on redis entry key
func (r *RedisRepository) playerSlotKey(slot string, playerID string) string {
	return "playerslot:" + slot + ":" + playerID
}

// ticketKey returns the key for the ticket
func (r *RedisRepository) ticketKey(ticketID string) string {
	return "ticket:" + ticketID
}
