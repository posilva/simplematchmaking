// Package lock provides the implementation of the locker output port using Redis.
package lock

import (
	"context"
	"errors"

	"github.com/redis/rueidis"
	redislock "github.com/redis/rueidis/rueidislock"
)

var (
	// ErrFailedToCreateLock is returned when the locker cannot be created
	ErrFailedToCreateLock = errors.New("failed to create locker")
)

// RedisLock is the Redis implementation of the Locker output port
type RedisLock struct {
	client rueidis.Client
	locker redislock.Locker
}

// NewRedisLock creates a new RedisLocker
func NewRedisLock(c rueidis.Client, keyMajority int32) (*RedisLock, error) {
	locker, err := redislock.NewLocker(redislock.LockerOption{
		ClientBuilder: func(co rueidis.ClientOption) (rueidis.Client, error) {
			return c, nil
		},
		KeyMajority:    keyMajority, // Use KeyMajority=1 if you have only one Redis instance. Also make sure that all your `Locker`s share the same KeyMajority.
		NoLoopTracking: true,        // Enable this to have better performance if all your Redis are >= 7.0.5.
	})

	if err != nil {
		return nil, errors.Join(ErrFailedToCreateLock, err)
	}
	return &RedisLock{
		client: c,
		locker: locker,
	}, nil
}

// Acquire acquires a lock
func (r *RedisLock) Acquire(ctx context.Context, key string) (context.Context, context.CancelFunc, error) {
	return r.locker.WithContext(ctx, key)
}
