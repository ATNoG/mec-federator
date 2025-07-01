package callbacks

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/router"
)

type FederationKduEnableCallback struct {
	services *router.Services
}

func NewEnableAppInstanceKDUCallback(services *router.Services) *FederationKduEnableCallback {
	return &FederationKduEnableCallback{
		services: services,
	}
}

func (f *FederationKduEnableCallback) HandleMessage(message *sarama.ConsumerMessage) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Handling KDU enable request: %v", msg)

	// TODO: Implement KDU enable logic

	// Send response to kafka
	_, err := f.services.KafkaClientService.Produce("responses", map[string]string{
		"msg_id": msg["msg_id"].(string),
		"status": "200",
	})
	if err != nil {
		log.Printf("Error sending response to kafka: %v", err)
		return
	}
}
