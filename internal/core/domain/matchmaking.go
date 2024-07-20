// Package domain implements the core domain models of the matchmaking service.
package domain

// TicketState represents the state of a ticket
type TicketState int

// TicketState enum
const (
	// TicketStateQueued state is the first state of a ticket
	TicketStateQueued TicketState = iota + 1
	// TicketStateMatched state is when a ticket has been matched with other tickets
	TicketStateMatched
	// TicketStateConfirmed state is when a ticket has been confirmed by the player
	TicketStateConfirmed
	// TicketStateCanceled state is when a ticket has been canceled by the player
	TicketStateCanceled
)

// Player represents a player in the matchmaking service
type Player struct {
	// ID is the player ID
	ID string `json:"id" msgpack:"id" mapstructure:"id"`
	// Ranking is the player ranking
	Ranking int `json:"ranking" msgpack:"ranking" mapstructure:"ranking"`
}

// Ticket represents a ticket in the matchmaking service
type Ticket struct {
	// ID is the ticket ID
	ID string `json:"id" msgpack:"id" mapstructure:"id"`
}

// Match represents a match in the matchmaking service
type Match struct {
	// ID is the match ID
	ID string `json:"id" msgpack:"id" mapstructure:"id"`
}

// MatchResult represents the result of a matchmaker operation
type MatchResult struct {
	// Match is the match data
	Match Match `json:"match" msgpack:"match" mapstructure:"match"`
	// Tickets are the tickets involved in the match
	Tickets []Ticket `json:"tickets" msgpack:"tickets" mapstructure:"tickets"`
}

// MatchmakerConfig represents the configuration of a matchmaker
type MatchmakerConfig struct {
	// Name is the matchmaker name
	Name string `json:"name" msgpack:"name" mapstructure:"name"`
	// MaxPlayers is the maximum number of players in a match
	MaxPlayers int `json:"maxPlayers" msgpack:"maxPlayers" mapstructure:"maxPlayers"`
	// Interval is the matchmaker interval
	Interval int `json:"interval" msgpack:"interval" mapstructure:"interval"`
	// Queue is the matchmaker queue
	Queue string `json:"queue" msgpack:"queue" mapstructure:"queue"`
}

// TicketStatus represents the data of a ticket saved in the repository
type TicketStatus struct {
	// ID is the ticket ID
	ID string `json:"id" msgpack:"id" mapstructure:"id"`
	// Timestamp is the state timestamp
	Timestamp int64 `json:"ts" msgpack:"ts" mapstructure:"ts"`
	// State is the ticket state
	State TicketState `json:"state" msgpack:"state" mapstructure:"state"`
	// PlayerID is the player ID
	PlayerID string `json:"uid" msgpack:"uid" mapstructure:"uid"`
}
