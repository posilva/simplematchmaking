// Package repository provides the repository implementation for the output port
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	logger ports.Logger
}

// NewRedisRepository creates a new RedisRepository
func NewRedisRepository(client rueidis.Client, codec ports.Codec, logger ports.Logger) *RedisRepository {
	return &RedisRepository{
		client: client,
		codec:  codec,
		logger: logger,
	}
}

// ReservePlayerSlot reserves a player slot in the queue
func (r *RedisRepository) ReservePlayerSlot(ctx context.Context, playerID string, slot string, ticketID string) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, redisCallTimeout)
	defer cancel()

	key := r.playerSlotKey(slot, playerID)
	value := "status:reserved:ticket:" + ticketID
	cmdSET := r.client.B().Set().Key(key).Value(value).Nx().ExSeconds(reservationTimeEx).Build()
	cmdGET := r.client.B().Get().Key(key).Build()
	resp := r.client.DoMulti(ctxWithTimeout, cmdSET, cmdGET)
	_, err := resp[0].AsBool()
	if err != nil {
		// if redis returns nil it means it already existed
		if !rueidis.IsRedisNil(err) {
			return "", errors.Join(ErrFailedToReservePlayerSlot, fmt.Errorf(" [SET] slot: '%v'  %v", slot, err))
		}
	}
	respGET, err := resp[1].AsBytes()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return "", errors.Join(ErrFailedToReservePlayerSlot, fmt.Errorf(" [GET] slot: '%v'  %v", slot, err))
		}
		return "", errors.Join(ErrSomethingOddHappened, fmt.Errorf("[GET] did returned 'nil' after SET: '%v'  %v", slot, err))

	}
	vGet := string(respGET)
	vGetTargetID := strings.Replace(vGet, "status:reserved:ticket:", "", -1)

	return vGetTargetID, nil
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
		return errors.Join(ErrFailedToEncodeTicket, err)
	}

	// TODO: this should have a TTL
	cmd := r.client.B().Set().Key(key).Value(string(value)).Build()
	err = r.client.Do(ctxWithTimeout, cmd).Error()
	if err != nil {
		if !rueidis.IsRedisNil(err) {
			return errors.Join(ErrFailedToUpdateTicket, err)
		}
		return errors.Join(ErrTicketNotFound, err)
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
			return domain.TicketStatus{}, errors.Join(ErrFailedToGetTicket, err)
		}
		return domain.TicketStatus{}, errors.Join(ErrTicketNotFound, err)
	}

	var status domain.TicketStatus
	err = r.codec.Decode([]byte(resp), &status)
	if err != nil {
		return domain.TicketStatus{}, errors.Join(ErrFailedToDecodeTicket, err)
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
			return domain.TicketStatus{}, errors.Join(ErrFailedToDeleteTicket, err)

		}
		return domain.TicketStatus{}, errors.Join(ErrTicketNotFound, err)
	}
	var status domain.TicketStatus
	err = r.codec.Decode(b, &status)
	if err != nil {
		return domain.TicketStatus{}, errors.Join(ErrFailedToDecodeTicket, err)
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
