package handler

// FindMatchInput is the input
type FindMatchInput struct {
	PlayerID string `json:"player_id"`
	Score    int    `json:"score"`
}

// FindMatchOutput is the output
type FindMatchOutput struct {
	TicketID string `json:"ticket_id"`
}

// GetMatchOutput is the output of getting a match
type GetMatchOutput struct {
	MatchID string   `json:"match_id"`
	Tickets []string `json:"tickets"`
}
