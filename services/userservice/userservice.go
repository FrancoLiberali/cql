package userservice

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/registry"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	validator "github.com/ditrit/badaas/validators"
)

// Create a new user
func NewUser(username, email, password string) (*models.User, error) {
	if !validator.ValidEmail(email) {
		return nil, fmt.Errorf("the provided email is not valid")
	}
	u := &models.User{
		Username: username,
		Email:    email,
		Password: basicauth.SaltAndHashPassword(password),
	}
	return u, registry.GetRegistry().UserRepository.Create(u)
}
