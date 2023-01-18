package userservice

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	validator "github.com/ditrit/badaas/validators"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UserService provide functions related to Users
type UserService interface {
	NewUser(username, email, password string) (*models.User, error)
	GetUser(dto.UserLoginDTO) (*models.User, httperrors.HTTPError)
}

// Check interface compliance
var _ UserService = (*userServiceImpl)(nil)

// The UserService concrete implementation
type userServiceImpl struct {
	userRepository repository.CRUDRepository[models.User, uuid.UUID]
	logger         *zap.Logger
}

// UserService constructor
func NewUserService(
	logger *zap.Logger,
	userRepository repository.CRUDRepository[models.User, uuid.UUID],
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

// Get user if the email and password provided are correct, return an error if not.
func (userService *userServiceImpl) GetUser(userLoginDTO dto.UserLoginDTO) (*models.User, httperrors.HTTPError) {
	users, herr := userService.userRepository.Find(squirrel.Eq{"email": userLoginDTO.Email}, nil, nil)
	if herr != nil {
		return nil, herr
	}
	if !users.HasContent {
		return nil, httperrors.NewErrorNotFound("user",
			fmt.Sprintf("no user found with email %q", userLoginDTO.Email))
	}

	user := users.Ressources[0]

	// Check password
	if !basicauth.CheckUserPassword(user.Password, userLoginDTO.Password) {
		return nil, httperrors.NewUnauthorizedError("wrong password", "the provided password is incorrect")
	}
	return user, nil
}
