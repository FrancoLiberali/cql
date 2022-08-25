package models

// Describe the current BADAAS instance
type BadaasServerInfo struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
