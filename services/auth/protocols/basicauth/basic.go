package basicauth

import (
	"golang.org/x/crypto/bcrypt"
)

// The bcrypt hashing cost (see [golang.org/x/crypto/bcrypt.Cost])
const cost = bcrypt.DefaultCost

// Salt and hash the password
func SaltAndHashPassword(password string) []byte {
	// bcrypt.GenerateFromPassword only return an error is the cost is not between 4 and 31, with cost=bcrypt.DefaultCost we are safe
	bytes, _ := bcrypt.GenerateFromPassword(
		[]byte(password),
		cost,
	)
	return bytes
}

// Check if the password is valid
func CheckUserPassword(passwordHash []byte, passwordToCheck string) bool {
	return bcrypt.CompareHashAndPassword(passwordHash, []byte(passwordToCheck)) == nil
}
