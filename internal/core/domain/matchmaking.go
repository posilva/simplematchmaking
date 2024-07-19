package domain

// Player ...
type Player struct {
	ID      string `json:"id"`
	Ranking int    `json:"ranking"`
}

// Ticket ...
type Ticket struct {
	ID string `json:"id"`
}

// Match ...
type Match struct {
	ID string `json:"id"`
}

// MatchResult ...
type MatchResult struct {
	Match   Match    `json:"match"`
	Tickets []Ticket `json:"tickets"`
}

// MatchmakerConfig ...
type MatchmakerConfig struct {
	MaxPlayers int    `json:"maxPlayers"`
	Name       string `json:"name"`
	Queue      string `json:"queue"`
}
