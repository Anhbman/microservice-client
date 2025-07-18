package config

type RabbitMQConfig struct {
	Host        string `json:"host" env:"RABBITMQ_HOST"`
	Port        int    `json:"port" env:"RABBITMQ_PORT"`
	Username    string `json:"username" env:"RABBITMQ_USERNAME"`
	Password    string `json:"password" env:"RABBITMQ_PASSWORD"`
	VirtualHost string `json:"virtual_host" env:"RABBITMQ_VHOST"`
	Exchange    string `json:"exchange" env:"RABBITMQ_EXCHANGE"`
}

func LoadRabbitMQConfig() (*RabbitMQConfig, error) {
	config := &RabbitMQConfig{
		Host:        GetEnv("RABBITMQ_HOST", "localhost"),
		Port:        GetEnvInt("RABBITMQ_PORT", 5672),
		Username:    GetEnv("RABBITMQ_USERNAME", "guest"),
		Password:    GetEnv("RABBITMQ_PASSWORD", "guest"),
		VirtualHost: GetEnv("RABBITMQ_VHOST", "/"),
		Exchange:    GetEnv("RABBITMQ_EXCHANGE", "client_events"),
	}

	return config, nil
}
