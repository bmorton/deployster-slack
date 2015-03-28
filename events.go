package main

type DeployEvent struct {
	Type          string `json:"type"`
	ServiceName   string `json:"service_name"`
	Version       string `json:"version"`
	Timestamp     string `json:"timestamp"`
	InstanceCount int    `json:"instance_count"`
}
