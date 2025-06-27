package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main1() {
	// Test connection
	fmt.Println("Testing RabbitMQ connection...")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected to RabbitMQ successfully")

	// Test channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	fmt.Println("✅ Channel opened successfully")

	// Test exchange declaration
	err = ch.ExchangeDeclare(
		"test_exchange", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare exchange:", err)
	}

	fmt.Println("✅ Exchange declared successfully")

	// Test queue declaration
	queue, err := ch.QueueDeclare(
		"test_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	fmt.Printf("✅ Queue declared successfully: %s\n", queue.Name)

	// Test queue binding
	err = ch.QueueBind(
		queue.Name,      // queue name
		"test.message",  // routing key
		"test_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to bind queue:", err)
	}

	fmt.Println("✅ Queue bound successfully")

	// Test publish message
	err = ch.Publish(
		"test_exchange", // exchange
		"test.message",  // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Hello RabbitMQ!"),
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		log.Fatal("Failed to publish message:", err)
	}

	fmt.Println("✅ Message published successfully")

	// Check queue stats
	inspectQueue, err := ch.QueueInspect("test_queue")
	if err != nil {
		log.Fatal("Failed to inspect queue:", err)
	}

	fmt.Printf("✅ Queue stats - Messages: %d, Consumers: %d\n",
		inspectQueue.Messages, inspectQueue.Consumers)

	fmt.Println("\n🎉 All RabbitMQ operations successful!")
}

// Alternative: Simple queue list check
func listQueues() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// This will create queue if it doesn't exist
	_, err = ch.QueueDeclare(
		"health_check", // name
		false,          // durable
		true,           // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	fmt.Println("Queue created successfully - RabbitMQ is working!")
}
