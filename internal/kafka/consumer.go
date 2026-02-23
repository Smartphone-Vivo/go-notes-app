package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"test-task/internal/domain"
	"test-task/internal/repository"
)

type Consumer struct {
	reader   *kafka.Reader
	noteRepo repository.NotesRepository
}

func NewConsumer(config *Config, noteRepo repository.NotesRepository) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		Topic:   config.Topic,
		GroupID: "notes-consumer-group",
	})

	return &Consumer{
		reader:   reader,
		noteRepo: noteRepo,
	}
}

func (c *Consumer) ConsumeMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			var note domain.Note
			if err := json.Unmarshal(msg.Value, &note); err != nil {
				continue
			}

			if err := c.noteRepo.CreateNote(ctx, note); err != nil {
				continue
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
