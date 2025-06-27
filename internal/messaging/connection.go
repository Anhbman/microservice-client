package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Config struct {
	Host        string
	Port        int
	Username    string
	Password    string
	VirtualHost string
}

func NewConnection(config Config) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.VirtualHost,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	fmt.Println("Connected to RabbitMQ")

	return conn, nil
}
