package userservice

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/services/auth/protocols/basicauth"
	"github.com/ditrit/badaas/utils/validators"
)

// UserService provide functions related to Users
type UserService interface {
	NewUser(username, email, password string) (*models.User, error)
	GetUser(dto.UserLoginDTO) (*models.User, error)
}

var ErrWrongPassword = errors.New("password is incorrect")

// Check interface compliance
var _ UserService = (*userServiceImpl)(nil)

// The UserService concrete implementation
type userServiceImpl struct {
	userRepository orm.CRUDRepository[models.User, model.UUID]
	logger         *zap.Logger
	db             *gorm.DB
}

// UserService constructor
func NewUserService(
	logger *zap.Logger,
	userRepository orm.CRUDRepository[models.User, model.UUID],
	db *gorm.DB,
) UserService {
	return &userServiceImpl{
		logger:         logger,
		userRepository: userRepository,
		db:             db,
	}
}

// Create a new user
func (userService *userServiceImpl) NewUser(username, email, password string) (*models.User, error) {
	sanitizedEmail, err := validators.ValidEmail(email)
	if err != nil {
		return nil, fmt.Errorf("the provided email is not valid")
	}

	u := &models.User{
		Username: username,
		Email:    sanitizedEmail,
		Password: basicauth.SaltAndHashPassword(password),
	}

	err = userService.userRepository.Create(userService.db, u)
	if err != nil {
		return nil, err
	}

	userService.logger.Info(
		"Successfully created a new user",
		zap.String("email", sanitizedEmail),
		zap.String("username", username),
	)

	return u, nil
}

// Get user if the email and password provided are correct, return an error if not.
func (userService *userServiceImpl) GetUser(userLoginDTO dto.UserLoginDTO) (*models.User, error) {
	user, err := userService.userRepository.QueryOne(
		userService.db,
		models.UserEmailCondition(orm.Eq(userLoginDTO.Email)),
	)
	if err != nil {
		return nil, err
	}

	// Check password
	if !basicauth.CheckUserPassword(user.Password, userLoginDTO.Password) {
		return nil, ErrWrongPassword
	}

	return user, nil
}
