package matchmaker

import (
	"context"

	"github.com/posilva/simplematchmaking/internal/core/domain"
	"github.com/redis/rueidis"
)

// RedisMatchmaker is the Matchmaker implementation using Redis
type RedisMatchmaker struct {
	client rueidis.Client
	config domain.MatchmakerConfig
}

// NewRedisMatchmaker creates a new RedisMatchmaker
func NewRedisMatchmaker(c rueidis.Client, cfg domain.MatchmakerConfig) (*RedisMatchmaker, error) {
	return &RedisMatchmaker{
		client: c,
		config: cfg,
	}, nil
}

// AddPlayer adds a player to the matchmaker
func (m *RedisMatchmaker) AddPlayer(ctx context.Context, p domain.Player) error {
	cmd := m.client.B().Zadd().Key("ranking:"+m.config.Queue).ScoreMember().ScoreMember(float64(p.Ranking), p.ID).Build()
	err := m.client.Do(ctx, cmd).Error()
	return err
}

// Match finds a match
func (m *RedisMatchmaker) Match(ctx context.Context) (domain.MatchResult, error) {
	return domain.MatchResult{
		Match: domain.Match{
			ID: "match1",
		},
		Tickets: []domain.Ticket{{ID: "ticket1"}, {ID: "ticket2"}},
	}, nil
}
