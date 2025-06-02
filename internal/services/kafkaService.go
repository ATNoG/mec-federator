package services

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/callbacks"
	"github.com/mankings/mec-federator/internal/config"
)

/*
 * KafkaService
 *	responsible for interacting with Kafka
 */

type KafkaServiceInterface interface {
	Produce(topic string, message interface{}) error
	StartConsumer(ctx context.Context, topic string, callback func(*sarama.ConsumerMessage)) error
}

type KafkaService struct {
	responseCallback      *callbacks.ResponseCallback
	newFederationCallback *callbacks.NewFederationCallback
}

func NewKafkaService() *KafkaService {
	kafkaServ := &KafkaService{
		responseCallback:      callbacks.NewResponseCallback(),
		newFederationCallback: callbacks.NewNewFederationCallback(),
	}

	// start response consumer
	err := kafkaServ.StartConsumer(context.Background(), "responses", kafkaServ.responseCallback.HandleMessage)
	if err != nil {
		log.Fatalf("Failed to start response consumer: %v", err)
	}

	// start new federation consumer
	err = kafkaServ.StartConsumer(context.Background(), "new_federation", kafkaServ.newFederationCallback.HandleMessage)
	if err != nil {
		log.Fatalf("Failed to start new federation consumer: %v", err)
	}

	return kafkaServ
}

// Produce a message to a topic in kafka
func (k *KafkaService) Produce(topic string, message interface{}) (string, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	// check if message has a msg_id field
	var jsonMessage map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonMessage); err != nil {
		return "", err
	}

	slog.Info("Producing message to topic", "topic", topic, "message", jsonMessage)

	if _, exists := jsonMessage["msg_id"]; !exists {
		jsonMessage["msg_id"] = uuid.New().String()
		// Re-marshal with the new msg_id
		bytes, err = json.Marshal(jsonMessage)
		if err != nil {
			return "", err
		}
	}

	slog.Info("Message after unmarshalling", "message", jsonMessage)

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = config.Producer.SendMessage(msg)
	if err != nil {
		return "", err
	}

	return jsonMessage["msg_id"].(string), nil
}

// StartConsumer starts consuming messages from a topic with the provided callback
func (k *KafkaService) StartConsumer(ctx context.Context, topic string, callback func(*sarama.ConsumerMessage)) error {
	consumer := config.Consumer

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		consumer.Close()
		return err
	}

	go func() {
		defer func() {
			partitionConsumer.Close()
			consumer.Close()
		}()

		for {
			select {
			case message := <-partitionConsumer.Messages():
				callback(message)
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming from %s topic: %v", topic, err)
			case <-ctx.Done():
				log.Printf("%s consumer shutting down", topic)
				return
			}
		}
	}()

	return nil
}
