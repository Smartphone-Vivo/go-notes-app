package kafka

import "os"

type Config struct {
	Brokers []string
	Topic   string
}

func NewConfig() *Config {
	return &Config{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   "notes-topic",
	}
}
