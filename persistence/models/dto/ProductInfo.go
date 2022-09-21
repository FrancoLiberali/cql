package dto

// Data Transfert Object Package

// Describe the Server Info payload
type DTOBadaasServerInfo struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
