package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/config"
)

/*
 * KafkaService
 *	responsible for interacting with Kafka
 */

type KafkaServiceInterface interface {
	Produce(topic string, message interface{}) error
	StartResponseConsumer(ctx context.Context) error
	GetResponse(msgID string) (map[string]interface{}, bool)
	WaitForResponse(msgID string, timeout time.Duration) (map[string]interface{}, error)
}

type KafkaService struct {
	responses map[string]map[string]interface{} // msgID -> response
	mu        sync.RWMutex                      // protects responses map
}

func NewKafkaService() *KafkaService {
	return &KafkaService{
		responses: make(map[string]map[string]interface{}),
	}
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

// StartResponseConsumer starts consuming messages from the "responses" topic
func (k *KafkaService) StartResponseConsumer(ctx context.Context) error {
	consumer := config.Consumer

	partitionConsumer, err := consumer.ConsumePartition("responses", 0, sarama.OffsetNewest)
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
				k.handleResponse(message)
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming from responses topic: %v", err)
			case <-ctx.Done():
				log.Println("Response consumer shutting down")
				return
			}
		}
	}()

	return nil
}

// handleResponse processes incoming response messages
func (k *KafkaService) handleResponse(message *sarama.ConsumerMessage) {
	var response map[string]interface{}
	if err := json.Unmarshal(message.Value, &response); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	msgID, exists := response["msg_id"]
	if !exists {
		log.Println("Response message missing msg_id field")
		return
	}

	msgIDStr, ok := msgID.(string)
	if !ok {
		log.Println("msg_id is not a string")
		return
	}

	k.mu.Lock()
	k.responses[msgIDStr] = response
	k.mu.Unlock()

	log.Printf("Stored response for msg_id: %s", msgIDStr)
}

// GetResponse retrieves a response by message ID (non-blocking)
func (k *KafkaService) GetResponse(msgID string) (map[string]interface{}, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	response, exists := k.responses[msgID]
	return response, exists
}

// WaitForResponse waits for a response with timeout (blocking)
func (k *KafkaService) WaitForResponse(msgID string, timeout time.Duration) (map[string]interface{}, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if response, exists := k.GetResponse(msgID); exists {
				return response, nil
			}
		case <-timeoutChan:
			return nil, ErrTimeout
		}
	}
}

// Custom errors
var (
	ErrTimeout = fmt.Errorf("timeout waiting for response")
)
