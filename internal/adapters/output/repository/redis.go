package repository

import (
	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/redis/rueidis"
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

// FindMatch finds a match given a player
func (r *RedisRepository) FindMatch(queue string, p domain.Player) (domain.Ticket, error) {
	_ = p
	_ = queue

	return domain.Ticket{
		ID: "ticket1",
	}, nil
}
