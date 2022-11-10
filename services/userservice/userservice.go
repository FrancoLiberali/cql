package userservice

import (
	"fmt"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	validator "github.com/ditrit/badaas/validators"
	"go.uber.org/zap"
)

// UserService provide functions related to Users
type UserService interface {
	NewUser(username, email, password string) (*models.User, error)
}

// Check interface compliance
var _ UserService = (*userServiceImpl)(nil)

// The UserService concrete implementation
type userServiceImpl struct {
	userRepository repository.CRUDRepository[models.User, uint]
	logger         *zap.Logger
}

// UserService constructor
func NewUserService(
	logger *zap.Logger,
	userRepository repository.CRUDRepository[models.User, uint],
) UserService {
	return &userServiceImpl{
		logger:         logger,
		userRepository: userRepository,
	}
}

// Create a new user
func (userService *userServiceImpl) NewUser(username, email, password string) (*models.User, error) {
	sanitizedEmail, err := validator.ValidEmail(email)
	if err != nil {
		return nil, fmt.Errorf("the provided email is not valid")
	}
	u := &models.User{
		Username: username,
		Email:    sanitizedEmail,
		Password: basicauth.SaltAndHashPassword(password),
	}
	httpError := userService.userRepository.Create(u)
	if httpError != nil {
		return nil, httpError
	}
	userService.logger.Info("Successfully created a new user",
		zap.String("email", sanitizedEmail), zap.String("username", username))

	return u, nil
}
