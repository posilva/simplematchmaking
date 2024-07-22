package queues

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/redis/rueidis"
	lock "github.com/redis/rueidis/rueidislock"
)

// RedisQueue is the Matchmaker implementation using Redis
type RedisQueue struct {
	client rueidis.Client
	name   string
	locker lock.Locker
}

// NewRedisQueue creates a new RedisQueue
func NewRedisQueue(c rueidis.Client, name string) *RedisQueue {
	locker, err := lock.NewLocker(lock.LockerOption{
		ClientBuilder: func(co rueidis.ClientOption) (rueidis.Client, error) {
			return c, nil
		},
		KeyMajority:    1,    // Use KeyMajority=1 if you have only one Redis instance. Also make sure that all your `Locker`s share the same KeyMajority.
		NoLoopTracking: true, // Enable this to have better performance if all your Redis are >= 7.0.5.
	})

	// This should not happend but nevertheless we should handle it
	// it will only happen if the client is not properly configured
	// when starting the queue
	if err != nil {
		panic(err)
	}
	return &RedisQueue{
		client: c,
		name:   name,
		locker: locker,
	}
}

// AddPlayer adds a player to the matchmaker
func (q *RedisQueue) AddPlayer(ctx context.Context, p domain.Player) error {
	cmd := q.client.B().Zadd().Key("ranking:"+q.name).ScoreMember().ScoreMember(float64(p.Ranking), p.ID).Build()
	err := q.client.Do(ctx, cmd).Error()
	return err
}

// Make finds a match
func (q *RedisQueue) Make(ctx context.Context, matchID string) (domain.MatchResult, error) {
	// acquire the lock "my_lock"
	ctxLock, cancel, err := q.locker.WithContext(ctx, q.name+":lock")
	if err != nil {
		return domain.MatchResult{}, err
	}
	defer cancel()

	cmd := q.client.B().Zrange().Key("ranking:" + q.name).Min("0").Max("-1").Rev().Build()
	res, err := q.client.Do(ctxLock, cmd).AsStrSlice()
	if err != nil {
		return domain.MatchResult{}, err
	}
	_ = res
	return domain.MatchResult{
		Match: domain.Match{
			ID: matchID,
		},
		Tickets: []domain.Ticket{{ID: "ticket1"}, {ID: "ticket2"}},
	}, nil
}

// Name returns the name of the queue
func (q *RedisQueue) Name() string {
	return q.name
}
