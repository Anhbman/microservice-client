package messaging

import (
	"context"
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

type EventHandler interface {
	Handle(ctx context.Context, event Event) error
	GetEventType() string
	GetQueueName() string
	GetRoutingKey() string
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	PublishEvent(ctx context.Context, event Event) error
	PublishEventWithRetry(ctx context.Context, event Event, maxRetries int) error
	RegisterEventType(eventType string, queueName string, routingKey string) error
	Close() error
}

// EventConsumer defines the interface for consuming events
type EventConsumer interface {
	Subscribe(handler EventHandler) error
	Start(ctx context.Context) error
	Stop() error
}

// EventRegistry holds event configuration
type EventRegistry struct {
	Events map[string]EventConfig `json:"events"`
}

type EventConfig struct {
	QueueName  string `json:"queue_name"`
	RoutingKey string `json:"routing_key"`
	Exchange   string `json:"exchange"`
	Durable    bool   `json:"durable"`
	AutoDelete bool   `json:"auto_delete"`
}

const (
	ClientExchange           = "client_events"
	RoutingKeyUserRegistered = "user_events.registered"
	// Event types
	EventTypeUserRegistered = "user_events.registered"
)
