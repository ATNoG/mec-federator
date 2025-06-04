package dto

import (
	"encoding/json"
)

// Net represents the network configuration
type Net map[string]string

// NodeSpec represents the specifications for a node
type NodeSpec struct {
	NumCPUCores  int     `json:"num_cpu_cores"`
	MemorySize   float64 `json:"memory_size"`
	AllocatedCPU int     `json:"allocated-cpu"`
	AllocatedMem int     `json:"allocated-mem"`
}

// Cluster represents a single cluster's data
type Cluster struct {
	Domain      string              `json:"domain"`
	Name        string              `json:"name"`
	K8sVersion  string              `json:"k8s-version"`
	VIMAccount  string              `json:"vim-account"`
	Description string              `json:"description"`
	Nets        Net                 `json:"nets"`
	NodeSpecs   map[string]NodeSpec `json:"nodeSpecs"`
}

// InfrastructureInfo represents the top-level structure of the Kafka message
type InfrastructureInfo struct {
	MsgID    string             `json:"msg_id"`
	Clusters map[string]Cluster `json:"-"` // Use `json:"-"` to exclude from direct unmarshaling
}

// UnmarshalJSON custom unmarshaler for InfrastructureInfo to handle dynamic cluster keys
func (km *InfrastructureInfo) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a temporary map to get all keys
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	km.Clusters = make(map[string]Cluster)

	for key, val := range raw {
		if key == "msg_id" {
			if err := json.Unmarshal(val, &km.MsgID); err != nil {
				return err
			}
		} else {
			// Assume any other key is a cluster ID
			var cluster Cluster
			if err := json.Unmarshal(val, &cluster); err != nil {
				return err
			}
			km.Clusters[key] = cluster
		}
	}
	return nil
}
