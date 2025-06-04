package callbacks

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/mankings/mec-federator/internal/models"
	"github.com/mankings/mec-federator/internal/models/dto"
)

// InfrastructureInfoCallback handles incoming infrastructure information messages from Kafka
type InfrastructureInfoCallback struct {
	latestMessage []byte
	latestZones   []models.ZoneDetails
	mu            sync.RWMutex
}

// NewInfrastructureInfoCallback creates a new InfrastructureInfoCallback instance
func NewInfrastructureInfoCallback() *InfrastructureInfoCallback {
	return &InfrastructureInfoCallback{}
}

// HandleMessage processes incoming cluster information messages
func (cc *InfrastructureInfoCallback) HandleMessage(message *sarama.ConsumerMessage) {
	cc.mu.Lock()
	cc.latestMessage = message.Value

	result, err := cc.UnmarshalInfrastructureInfo()
	if err != nil {
		log.Println("Error unmarshalling infrastructure info:", err)
		return
	}

	cc.latestZones = result

	cc.mu.Unlock()
}

// GetLatestZones returns the latest zones
func (cc *InfrastructureInfoCallback) GetLatestZones() []models.ZoneDetails {
	cc.mu.RLock()
	defer cc.mu.RUnlock()
	return cc.latestZones
}

// UnmarshalInfrastructureInfo unmarshals the latest cluster information message
func (cc *InfrastructureInfoCallback) UnmarshalInfrastructureInfo() ([]models.ZoneDetails, error) {
	if cc.latestMessage == nil {
		return nil, nil
	}

	var infraInfo dto.InfrastructureInfo
	if err := infraInfo.UnmarshalJSON(cc.latestMessage); err != nil {
		return nil, err
	}

	// extract the zones from the message
	availableZones := make([]models.ZoneDetails, 0)

	// The message format is {cluster1Id: {clusterInfo}, cluster2Id: {clusterInfo}, ...}
	for clusterId, clusterInfo := range infraInfo.Clusters {
		zone := models.ZoneDetails{
			ZoneId: clusterId,
			VimId:  clusterInfo.VIMAccount,
		}

		availableZones = append(availableZones, zone)
	}

	return availableZones, nil
}
