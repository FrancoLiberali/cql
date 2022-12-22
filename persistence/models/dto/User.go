package dto

// Data Transfert Object Package

// Login DTO
type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
