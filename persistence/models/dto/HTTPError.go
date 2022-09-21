package dto

// Data Transfert Object Package

// Describe the HTTP Error payload
type DTOHTTPError struct {
	Error   string `json:"err"`
	Message string `json:"msg"`
	Status  string `json:"status"`
}
