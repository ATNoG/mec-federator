package callbacks

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/router"
	"github.com/mankings/mec-federator/internal/utils"
)

type FederationUpdateAppCallback struct {
	services *router.Services
}

func NewFederationUpdateAppCallback(services *router.Services) *FederationUpdateAppCallback {
	return &FederationUpdateAppCallback{
		services: services,
	}
}

func (f *FederationUpdateAppCallback) HandleMessage(message *sarama.ConsumerMessage) {
	utils.TimeCallback("FederationUpdateAppCallback.HandleMessage", func() {
		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-init-update-app",
			Message: "reset",
			Value:   nil,
		})

		log.Printf("Received update app message from topic %s, partition %d, offset %d",
			message.Topic, message.Partition, message.Offset)

		var msg map[string]interface{}
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return
		}

		log.Printf("Processing update app request with message ID: %s", msg["msg_id"])

		msgId := msg["msg_id"].(string)
		f.handleUpdateApp(msgId, msg)

		utils.SendResultsMessage(utils.ResultsMessage{
			Name:    "oo-done-update-app",
			Message: "",
			Value:   nil,
		})
	})
}

func (f *FederationUpdateAppCallback) handleUpdateApp(msgId string, msg map[string]interface{}) {

}
