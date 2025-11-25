package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResultsMessage struct {
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

var resultsEndpoint = "http://10.255.42.59:8000/alert"

func SendResultsMessage(message ResultsMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("results - failed to marshal results message: %w", err)
	}

	resp, err := http.Post(resultsEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("results - failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("results - received non-success status code: %d", resp.StatusCode)
	}

	return nil
}
