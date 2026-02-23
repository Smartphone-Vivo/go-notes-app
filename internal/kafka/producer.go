package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"test-task/internal/domain"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(config *Config) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(config.Brokers...),
		Topic:                  config.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) PublishNote(ctx context.Context, note domain.Note) error {

	data, err := json.Marshal(note)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(note.ID),
		Value: data,
		Time:  time.Now(),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
