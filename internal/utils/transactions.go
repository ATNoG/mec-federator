package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
)

func GenerateIdempotencyKey(request interface{}) string {
	hash := sha256.New()
	json, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return ""
	}
	hash.Write(json)
	return hex.EncodeToString(hash.Sum(nil))
}
