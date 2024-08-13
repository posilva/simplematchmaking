package health

import (
	"context"
	"time"

	"github.com/redis/rueidis"
)

type RedisCheck struct {
	client rueidis.Client
}

func NewRedisCheck(client rueidis.Client) *RedisCheck {
	return &RedisCheck{client: client}
}

func (r *RedisCheck) Pass() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := r.client.B().Ping().Build()
	return nil == r.client.Do(ctx, cmd).Error()
}

func (r *RedisCheck) Name() string {
	return "redis"
}
