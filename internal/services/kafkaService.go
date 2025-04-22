package services

import (
	"context"

	"github.com/segmentio/kafka-go"
)

/*
 * KafkaService
 *	responsible for interacting with Kafka
 */

type KafkaServiceInterface struct {
}

type KafkaService struct {
	writer *kafka.Writer
}

func NewKafkaService(writer *kafka.Writer) *KafkaService {
	return &KafkaService{
		writer: writer,
	}
}

func (k *KafkaService) Produce(ctx context.Context, msg kafka.Message) error {
	return k.writer.WriteMessages(ctx, msg)
}
