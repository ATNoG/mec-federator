package callbacks

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// ResponseCallback handles incoming response messages from Kafka
type ResponseCallback struct {
	responses  map[string]map[string]interface{} // msgID -> response
	timestamps map[string]time.Time              // msgID -> timestamp
	mu         sync.RWMutex                      // protects responses map
}

// NewResponseCallback creates a new ResponseCallback instance
func NewResponseCallback() *ResponseCallback {
	return &ResponseCallback{
		responses:  make(map[string]map[string]interface{}),
		timestamps: make(map[string]time.Time),
	}
}

// HandleMessage processes incoming response messages
func (rc *ResponseCallback) HandleMessage(message *sarama.ConsumerMessage) {
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

	rc.mu.Lock()
	rc.responses[msgIDStr] = response
	rc.timestamps[msgIDStr] = time.Now()
	rc.mu.Unlock()

	log.Printf("Stored response for msg_id: %s", msgIDStr)
}

// GetResponse retrieves a response by message ID (non-blocking)
func (rc *ResponseCallback) GetResponse(msgID string) (map[string]interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	response, exists := rc.responses[msgID]
	return response, exists
}

// WaitForResponse waits for a response with timeout (blocking)
func (rc *ResponseCallback) WaitForResponse(msgID string, timeout time.Duration) (map[string]interface{}, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if response, exists := rc.GetResponse(msgID); exists {
				return response, nil
			}
		case <-timeoutChan:
			return nil, errors.New("timeout waiting for response")
		}
	}
}

// CleanupOldMessages removes messages older than the specified TTL
func (rc *ResponseCallback) CleanupOldMessages(messageTTL time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	now := time.Now()
	for msgID, timestamp := range rc.timestamps {
		if now.Sub(timestamp) > messageTTL {
			delete(rc.responses, msgID)
			delete(rc.timestamps, msgID)
			log.Printf("Cleaned up old message: %s", msgID)
		}
	}
}
