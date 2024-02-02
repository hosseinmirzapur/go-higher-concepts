package kafka

import "os"

type Config struct {
	Host string
	Port string
}

func KafkaConfig() *Config {
	return &Config{
		Host: os.Getenv("KAFKA_HOST"),
		Port: os.Getenv(os.Getenv("KAFKA_PORT")),
	}
}
