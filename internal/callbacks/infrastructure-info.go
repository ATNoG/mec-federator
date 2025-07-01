package callbacks

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
	"github.com/mankings/mec-federator/internal/router"
)

// InfrastructureInfoCallback handles incoming infrastructure information messages from Kafka
type InfrastructureInfoCallback struct {
	latestMessage []byte
	latestZones   []models.ZoneDetails
	mu            sync.RWMutex

	services *router.Services
}

// NewInfrastructureInfoCallback creates a new InfrastructureInfoCallback instance
func NewInfrastructureInfoCallback(services *router.Services) *InfrastructureInfoCallback {
	return &InfrastructureInfoCallback{
		services: services,
	}
}

// HandleMessage processes incoming cluster information messages
func (cc *InfrastructureInfoCallback) HandleMessage(message *sarama.ConsumerMessage) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// Store the latest message
	cc.latestMessage = message.Value

	// Unmarshal the message directly
	var infraInfo dto.InfrastructureInfo
	if err := infraInfo.UnmarshalJSON(message.Value); err != nil {
		log.Println("Error unmarshalling infrastructure info:", err)
		return
	}

	// Extract the zones from the message
	availableZones := make([]models.ZoneDetails, 0)

	// The message format is {cluster1Id: {clusterInfo}, cluster2Id: {clusterInfo}, ...}
	for clusterId, clusterInfo := range infraInfo.Clusters {
		zone := models.ZoneDetails{
			ZoneId: clusterId,
			VimId:  clusterInfo.VIMAccount,
		}
		availableZones = append(availableZones, zone)
	}

	// Store the latest zones
	cc.latestZones = availableZones

	// Update the zones in the database
	if err := cc.services.ZoneService.UpdateLocalZones(availableZones); err != nil {
		log.Println("Error updating local zones:", err)
		return
	}

}
