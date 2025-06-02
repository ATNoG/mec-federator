package callbacks

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

// NewFederationCallback handles incoming new federation messages from Kafka
type NewFederationCallback struct {
}

// NewNewFederationCallback creates a new NewFederationCallback instance
func NewNewFederationCallback() *NewFederationCallback {
	return &NewFederationCallback{}
}

// HandleMessage processes incoming new federation messages
func (nc *NewFederationCallback) HandleMessage(message *sarama.ConsumerMessage) {
	var response map[string]interface{}
	if err := json.Unmarshal(message.Value, &response); err != nil {
		log.Printf("Error unmarshaling response message: %v", err)
		return
	}

	federationEndpoint := response["federation_endpoint"].(string)
	authEndpoint := response["auth_endpoint"].(string)
	clientId := response["client_id"].(string)
	clientSecret := response["client_secret"].(string)

	msgID, ok := response["msg_id"].(string)
	if !ok {
		log.Println("msg_id is not a string")
		return
	}

	log.Printf("New federation created: %s", msgID)
	log.Printf("Federation endpoint: %s", federationEndpoint)
	log.Printf("Auth endpoint: %s", authEndpoint)
	log.Printf("Client ID: %s", clientId)
	log.Printf("Client secret: %s", clientSecret)

	// send response after
}
