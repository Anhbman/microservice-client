package messaging

import (
	"client/internal/config"
	"fmt"

	"github.com/streadway/amqp"
)

// type Config struct {
// 	Host        string
// 	Port        int
// 	Username    string
// 	Password    string
// 	VirtualHost string
// }

func NewConnection(conf config.RabbitMQConfig) (*amqp.Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.VirtualHost,
	)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	return conn, nil
}
