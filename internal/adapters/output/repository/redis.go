// Package repository provides the repository implementation for the output port
package repository

import (
	"context"
	"time"

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
}

// NewRedisRepository creates a new RedisRepository
func NewRedisRepository(client rueidis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
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
