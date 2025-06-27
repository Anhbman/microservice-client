package messaging

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
	Source    string          `json:"source"`
}

type UserRegisterPayload struct {
	Name 	 string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	ClientExchange = "client_events"
	RoutingKeyUserRegistered = "user_events.registered"
	// Event types
	EventTypeUserRegistered = "user_events.registered"
)
