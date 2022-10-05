package validator

import "net/mail"

// Check if the email is a valid email (according to [RFC 5322])
func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
