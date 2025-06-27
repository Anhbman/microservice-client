package messaging

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	fmt.Println("✅ Channel opened successfully")

	err = ch.ExchangeDeclare(
		ClientExchange, // name
		"topic",        // type
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
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *Publisher) PublishEvent(routingKey, queue string, eventType string, payload interface{}) error {
	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Tạo event
	event := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Payload:   payloadBytes,
		Timestamp: time.Now(),
		Source:    "client-service",
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Đảm bảo queue và binding tồn tại
	if _, err := p.channel.QueueDeclare(
		queue, true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	if err := p.channel.QueueBind(
		queue, routingKey, ClientExchange, false, nil,
	); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	// Publish
	return p.channel.Publish(
		ClientExchange, routingKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventBytes,
			Timestamp:   time.Now(),
		},
	)
}

func (p *Publisher) PublishUserRegisteredEvent(payload UserRegisterPayload) error {
	// payloadBytes, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal payload: %w", err)
	// }

	// err = p.DeclareAndBindQueue()
	// if err != nil {
	// 	// handle error
	// }

	// event := Event{
	// 	ID:        uuid.New().String(),
	// 	Type:      EventTypeUserRegistered,
	// 	Payload:   payloadBytes,
	// 	Timestamp: time.Now(),
	// 	Source:    "client-service",
	// }

	// eventBytes, err := json.Marshal(event)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal event: %w", err)
	// }

	err := p.PublishEvent(
		RoutingKeyUserRegistered,
		"user_registered_queue",
		EventTypeUserRegistered,
		payload,
	)
	if err != nil {
		return fmt.Errorf("failed to publish user registered event: %w", err)
	}
	return nil
}

// func (p *Publisher) DeclareAndBindQueue() error {
// 	queueName := "user_registered_queue"
// 	routingKey := "user.registered"
// 	exchange := "user_events"

// 	// Khai báo queue
// 	_, err := p.channel.QueueDeclare(
// 		queueName, // tên queue
// 		true,      // durable
// 		false,     // auto-delete
// 		false,     // exclusive
// 		false,     // no-wait
// 		nil,       // arguments
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to declare queue: %w", err)
// 	}

// 	// Bind queue vào exchange với routing key
// 	err = p.channel.QueueBind(
// 		queueName,  // queue name
// 		routingKey, // routing key
// 		exchange,   // exchange
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to bind queue: %w", err)
// 	}

// 	return nil
// }

func (p *Publisher) Close() error {
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
