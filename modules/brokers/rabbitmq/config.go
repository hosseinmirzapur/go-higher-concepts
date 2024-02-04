package rabbitmq

import (
	"os"
	"time"
)

type Config struct {
	URI     string
	Timeout time.Duration
}

func RabbitMQConfig() *Config {
	return &Config{
		URI:     os.Getenv("RABBITMQ_URI"),
		Timeout: 5 * time.Second,
	}
}
