package dto

// Data Transfer Object Package

// Describe the HTTP Error payload
type HTTPError struct {
	Error   string `json:"err"`
	Message string `json:"msg"`
	Status  string `json:"status"`
}
