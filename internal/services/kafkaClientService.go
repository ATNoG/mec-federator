package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mankings/mec-federator/internal/config"
)

type consumerGroupHandler struct {
	topic    string
	callback func(*sarama.ConsumerMessage)
	logs     bool
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}
			if h.logs {
				log.Printf("Received message from topic '%s': %s", h.topic, string(message.Value))
			}
			h.callback(message)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

/*
 * KafkaService
 *	responsible for interacting with Kafka
 */

type KafkaClientServiceInterface interface {
	Produce(topic string, message interface{}) error
	StartConsumer(ctx context.Context, topic string, callback func(*sarama.ConsumerMessage)) error
	SendResponse(msgId, status, message string) error
}

type KafkaClientService struct {
	callbacks  map[string]func(*sarama.ConsumerMessage)
	responses  map[string]map[string]interface{} // msgID -> response
	timestamps map[string]time.Time              // msgID -> timestamp
	mu         sync.RWMutex                      // protects responses map
}

func NewKafkaClientService() *KafkaClientService {
	return &KafkaClientService{
		callbacks:  make(map[string]func(*sarama.ConsumerMessage)),
		responses:  make(map[string]map[string]interface{}),
		timestamps: make(map[string]time.Time),
	}
}

// Produce a message to a topic in kafka
func (k *KafkaClientService) Produce(topic string, message interface{}) (string, error) {
	bytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	// check if message has a msg_id field
	var jsonMessage map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonMessage); err != nil {
		return "", err
	}

	if _, exists := jsonMessage["msg_id"]; !exists {
		jsonMessage["msg_id"] = uuid.New().String()
		// Re-marshal with the new msg_id
		bytes, err = json.Marshal(jsonMessage)
		if err != nil {
			return "", err
		}
	}

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
func (k *KafkaClientService) StartConsumer(ctx context.Context, topic string, callback func(*sarama.ConsumerMessage), logs bool) error {
	brokers := []string{config.AppConfig.KafkaHost + ":" + config.AppConfig.KafkaPort}
	consumerConfig := sarama.NewConfig()
	consumerConfig.Net.SASL.Enable = true
	consumerConfig.Net.SASL.User = config.AppConfig.KafkaUsername
	consumerConfig.Net.SASL.Password = config.AppConfig.KafkaPassword
	
	// Consumer group configuration
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerConfig.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(brokers, "mec-federator-group", consumerConfig)
	if err != nil {
		return err
	}

	log.Println("Starting consumer for topic", topic)

	k.callbacks[topic] = callback

	// If this is the responses topic, wrap the callback with response handling
	if topic == "responses" {
		originalCallback := callback
		callback = func(message *sarama.ConsumerMessage) {
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
			k.timestamps[msgIDStr] = time.Now()
			k.mu.Unlock()

			// Call original callback if provided
			if originalCallback != nil {
				originalCallback(message)
			}
		}
	}

	// Create consumer group handler
	handler := &consumerGroupHandler{
		topic:    topic,
		callback: callback,
		logs:     logs,
	}

	go func() {
		defer consumerGroup.Close()
		
		for {
			select {
			case <-ctx.Done():
				log.Printf("%s consumer shutting down", topic)
				return
			default:
				if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
					log.Printf("Error consuming from %s topic: %v", topic, err)
				}
			}
		}
	}()

	log.Println("Consumer started for topic", topic)

	return nil
}

// Add these methods to KafkaClientService
func (k *KafkaClientService) GetResponse(msgID string) (map[string]interface{}, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	response, exists := k.responses[msgID]
	return response, exists
}

func (k *KafkaClientService) WaitForResponse(msgID string, timeout time.Duration) (map[string]interface{}, error) {
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
			return nil, errors.New("timeout waiting for response")
		}
	}
}

func (k *KafkaClientService) CleanupOldMessages(messageTTL time.Duration) {
	k.mu.Lock()
	defer k.mu.Unlock()

	now := time.Now()
	for msgID, timestamp := range k.timestamps {
		if now.Sub(timestamp) > messageTTL {
			delete(k.responses, msgID)
			delete(k.timestamps, msgID)
		}
	}
}

// SendResponse sends a standardized response message to the responses topic
func (k *KafkaClientService) SendResponse(msgId, status, message string) error {
	response := map[string]string{
		"msg_id": msgId,
		"status": status,
	}

	// Add message field if it's not empty
	if message != "" {
		response["message"] = message
	}

	_, err := k.Produce("responses", response)
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return err
	}

	return nil
}
