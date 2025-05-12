package services

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/config"
)

/*
 * KafkaService
 *	responsible for interacting with Kafka
 */

type KafkaServiceInterface struct {
}

type KafkaService struct {
	topics []string
}

func NewKafkaService() *KafkaService {
	return &KafkaService{
		topics: []string{"federator-topic"},
	}
}

// Produce a message to a topic in kafka
func (k *KafkaService) Produce(topic string, message interface{}) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = config.Producer.SendMessage(msg)
	return err
}
