package matchmaker

import "github.com/posilva/simplematchmaking/internal/core/domain"

// RankingQueue is the queue that holds the ranking
type RankingQueue struct {
	score int
}

// MemoryMatchmaker is the Matchmaker implementation using memory
type MemoryMatchmaker struct {
	config domain.MatchmakerConfig
	queue  map[string]RankingQueue
}

// NewMemoryMatchmaker creates a new MemoryMatchmaker
func NewMemoryMatchmaker(cfg domain.MatchmakerConfig) *MemoryMatchmaker {
	return &MemoryMatchmaker{
		config: cfg,
		queue:  make(map[string]RankingQueue),
	}
}

// Match finds a match
func (m *MemoryMatchmaker) Match() (domain.MatchResult, error) {
	return domain.MatchResult{
		Match: domain.Match{
			ID: "match1",
		},
		Tickets: []domain.Ticket{{ID: "ticket1"}, {ID: "ticket2"}},
	}, nil
}
