package messaging

import (
	"client/internal/config"
	"context"
	"fmt"
	"sync"
)

type EventManager struct {
	publisher EventPublisher
	consumers map[string]EventConsumer
	config    config.RabbitMQConfig
	mutex     sync.RWMutex
}

func NewEventManager(config config.RabbitMQConfig) (*EventManager, error) {
	conn, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	publisher, err := NewPublisher(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &EventManager{
		publisher: publisher,
		consumers: make(map[string]EventConsumer),
		config:    config,
	}, nil
}

func (em *EventManager) RegisterEvents() error {
	// Register all event types
	events := map[string]struct {
		queue      string
		routingKey string
	}{
		UserRegisteredEvent: {"user_registered_queue", "user.registered"},
		UserUpdatedEvent:    {"user_updated_queue", "user.updated"},
		UserDeletedEvent:    {"user_deleted_queue", "user.deleted"},
		CakeCreatedEvent:    {"cake_created_queue", "cake.created"},
		CakeUpdatedEvent:    {"cake_updated_queue", "cake.updated"},
		CakeDeletedEvent:    {"cake_deleted_queue", "cake.deleted"},
		OrderCreatedEvent:   {"order_created_queue", "order.created"},
		OrderUpdatedEvent:   {"order_updated_queue", "order.updated"},
	}

	for eventType, config := range events {
		if err := em.publisher.RegisterEventType(eventType, config.queue, config.routingKey); err != nil {
			return fmt.Errorf("failed to register event type %s: %w", eventType, err)
		}
	}

	return nil
}

func (em *EventManager) PublishEvent(ctx context.Context, eventType string, payload interface{}) error {
	event := NewEventBuilder().
		WithType(eventType).
		WithPayload(payload).
		Build()

	return em.publisher.PublishEventWithRetry(ctx, event, 3)
}

func (em *EventManager) GetPublisher() EventPublisher {
	return em.publisher
}

func (em *EventManager) Close() error {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	// Close all consumers
	for _, consumer := range em.consumers {
		if err := consumer.Stop(); err != nil {
		}
	}

	// Close publisher
	return em.publisher.Close()
}
