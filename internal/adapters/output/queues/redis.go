package queues

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/redis/rueidis"
)

// RedisQueue is the Matchmaker implementation using Redis
type RedisQueue struct {
	client rueidis.Client
	name   string
}

// NewRedisQueue creates a new RedisQueue
func NewRedisQueue(c rueidis.Client, name string) *RedisQueue {
	return &RedisQueue{
		client: c,
		name:   name,
	}
}

// AddPlayer adds a player to the matchmaker
func (q *RedisQueue) AddPlayer(ctx context.Context, p domain.Player) error {
	cmd := q.client.B().Zadd().Key("ranking:"+q.name).ScoreMember().ScoreMember(float64(p.Ranking), p.ID).Build()
	err := q.client.Do(ctx, cmd).Error()
	return err
}
