package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	broker string
}

func NewProducer(broker string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{
		broker: broker,
		writer: writer,
	}
}

func (p *Producer) SendEvent(topic string, event interface{}) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   nil,
		Value: jsonData,
		Topic: topic,
	}

	err = p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write message: %v", err)
		return err
	}
	log.Printf("Produced event: %s", string(jsonData))
	return nil
}

func (p *Producer) Close() {
	p.writer.Close()
}
