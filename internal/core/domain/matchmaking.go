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
