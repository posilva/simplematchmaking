// Package repository provides the repository implementation for the output port
package repository

import (
	"context"
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
	key := r.playerSlotKey(slot, playerID)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()
	value := "status:reserved:ticket:" + ticketID
	cmd := r.client.B().Set().Key(key).Value(value).Nx().ExSeconds(reservationTimeEx).Build()
	resp, err := r.client.Do(ctxWithTimeout, cmd).AsBool()
	if err != nil {
		return false, err
	}
	return resp, nil
}

// playerSlotKey returns the key for the player slot to use on redis entry key
func (r *RedisRepository) playerSlotKey(slot string, playerID string) string {
	return "playerslot:" + slot + ":" + playerID
}

// ticketKey returns the key for the ticket
func (r *RedisRepository) ticketKey(ticketID string) string {
	return "ticket:" + ticketID
}

// UpdateTicketStatus updates the ticket status
func (r *RedisRepository) UpdateTicketStatus(ctx context.Context, status domain.TicketStatus) error {
	key := r.ticketKey(status.ID)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	value, err := r.codec.Encode(status)
	if err != nil {
		return err
	}

	// TODO: this should have a TTL
	cmd := r.client.B().Set().Key(key).Value(string(value)).Build()
	err = r.client.Do(ctxWithTimeout, cmd).Error()
	if err != nil {
		return err
	}
	return nil
}

// GetTicketStatus gets the ticket status
func (r *RedisRepository) GetTicketStatus(ctx context.Context, ticketID string) (domain.TicketStatus, error) {
	key := r.ticketKey(ticketID)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	cmd := r.client.B().Get().Key(key).Build()
	resp, err := r.client.Do(ctxWithTimeout, cmd).AsBytes()
	if err != nil {
		return domain.TicketStatus{}, err
	}
	var status domain.TicketStatus
	err = r.codec.Decode([]byte(resp), &status)
	if err != nil {
		return domain.TicketStatus{}, err
	}
	return status, nil
}
