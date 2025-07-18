package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	registry *EventRegistry
	mutex    sync.RWMutex
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	fmt.Println("✅ Channel opened successfully")

	err = ch.ExchangeDeclare(
		ClientExchange, // name
		Topic,          // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	fmt.Println("✅ Exchange declared successfully")

	return &Publisher{
		conn:     conn,
		channel:  ch,
		registry: &EventRegistry{Events: make(map[string]EventConfig)},
	}, nil
}

func (p *Publisher) RegisterEventType(eventType string, queueName string, routingKey string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	config := EventConfig{
		QueueName:  queueName,
		RoutingKey: routingKey,
		Exchange:   ClientExchange,
		Durable:    true,
		AutoDelete: false,
	}

	// Declare queue and binding
	if err := p.setupQueueAndBinding(config); err != nil {
		return fmt.Errorf("failed to setup queue and binding: %w", err)
	}

	p.registry.Events[eventType] = config
	log.Infof("Event type %s registered with queue %s and routing key %s", eventType, queueName, routingKey)
	return nil
}

func (p *Publisher) setupQueueAndBinding(config EventConfig) error {
	// Declare queue
	_, err := p.channel.QueueDeclare(
		config.QueueName,
		config.Durable,
		config.AutoDelete,
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", config.QueueName, err)
	}

	// Bind queue to exchange
	err = p.channel.QueueBind(
		config.QueueName,
		config.RoutingKey,
		config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue %s: %w", config.QueueName, err)
	}

	return nil
}

func (p *Publisher) PublishEvent(ctx context.Context, event Event) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Get event configuration
	eventConfig, exists := p.registry.Events[event.Type]
	if !exists {
		return fmt.Errorf("event type %s not registered", event.Type)
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Ensure connection is alive
	if p.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}

	// Set event metadata
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Source == "" {
		event.Source = "client-service"
	}

	// Marshal event
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish message
	err = p.channel.Publish(
		eventConfig.Exchange,
		eventConfig.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBytes,
			Timestamp:   time.Now(),
			MessageId:   event.ID,
		},
	)

	if err != nil {
		log.Errorf("failed to publish event %s: %v", event.Type, err)
		return fmt.Errorf("failed to publish event: %w", err)
	}
	log.Infof("Event %s published successfully", event.Type)
	return nil
}

func (p *Publisher) PublishEventWithRetry(ctx context.Context, event Event, maxRetries int) error {
	var err error
	for i := 0; i <= maxRetries; i++ {
		err = p.PublishEvent(ctx, event)
		if err == nil {
			return nil
		}

		if i < maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second * time.Duration(i+1)):
				// Exponential backoff
			}
		}
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

func (p *Publisher) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var err1, err2 error
	if p.channel != nil {
		err1 = p.channel.Close()
	}
	if p.conn != nil {
		err2 = p.conn.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}
