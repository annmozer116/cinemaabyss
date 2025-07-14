package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cinemaabyss/microservices/events/model"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker string, topics []string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		GroupTopics: topics,
		GroupID:     "events-group",
		MinBytes:    10e1,
		MaxBytes:    10e6,
	})
	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) ConsumeEvents(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error consuming message: %v", err)
			break
		}

		var event model.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		log.Printf("Consumed event: %s with data: %s", event.Type, event.Data)
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}
