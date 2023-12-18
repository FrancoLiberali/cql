package dto

// LoginSuccess is a dto returned to the client when the authentication is successful.
type LoginSuccess struct {
	Email    string `json:"email"`
	ID       string `json:"id"`
	Username string `json:"username"`
}
