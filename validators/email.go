package validator

import "net/mail"

// Check if the email is a valid email (according to [RFC 5322])
func ValidEmail(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}

	return addr.Address, nil
}
