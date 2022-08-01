package userservice_test

import (
	"testing"

	"github.com/ditrit/badaas/httperrors"
	repositorymocks "github.com/ditrit/badaas/mocks/persistence/repository"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewUserService(t *testing.T) {
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userRespositoryMock.On("Create", mock.Anything).Return(nil)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	user, err := userService.NewUser("bob", "bob@email.com", "1234")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "bob", user.Username)
	assert.Equal(t, "bob@email.com", user.Email)
	assert.NotEqual(t, "1234", user.Password)

	// Checking logs
	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Successfully created a new user", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "email", Type: zapcore.StringType, String: "bob@email.com"},
		{Key: "username", Type: zapcore.StringType, String: "bob"},
	}, log.Context)
	assert.Equal(t, zap.InfoLevel, log.Level)
}

func TestNewUserServiceDatabaseError(t *testing.T) {
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userRespositoryMock.On(
		"Create", mock.Anything,
	).Return(
		httperrors.NewInternalServerError("database error", "test error", nil),
	)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	user, err := userService.NewUser("bob", "bob@email.com", "1234")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Checking logs
	assert.Equal(t, 0, observedLogs.Len())
}

func TestNewUserServiceEmailNotValid(t *testing.T) {
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)

	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	user, err := userService.NewUser("bob", "bob@", "1234")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Checking logs
	assert.Equal(t, 0, observedLogs.Len())
}

func TestGetUser(t *testing.T) {
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	userRespositoryMock.On(
		"Create", mock.Anything,
	).Return(
		nil,
	)
	user, err := userService.NewUser("bob", "bob@email.com", "1234")

	require.NoError(t, err)
	userRespositoryMock.On(
		"Find", mock.Anything, nil, nil,
	).Return(
		pagination.NewPage([]*models.User{user}, 1, 10, 50),
		nil,
	)

	userFound, err := userService.GetUser(dto.UserLoginDTO{Email: "bob@email.com", Password: "1234"})
	require.NoError(t, err)
	assert.Equal(t, user, userFound)
	// Checking logs
	assert.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Successfully created a new user", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "email", Type: zapcore.StringType, String: "bob@email.com"},
		{Key: "username", Type: zapcore.StringType, String: "bob"},
	}, log.Context)
	assert.Equal(t, zap.InfoLevel, log.Level)
}

func TestGetUserNoUserFound(t *testing.T) {
	// creating logger
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	userRespositoryMock.On(
		"Find", mock.Anything, nil, nil,
	).Return(
		nil,
		httperrors.NewErrorNotFound("user", "user with email bobnotfound@email.com"),
	)

	userFound, err := userService.GetUser(dto.UserLoginDTO{Email: "bobnotfound@email.com", Password: "1234"})
	require.Error(t, err)
	assert.Nil(t, userFound)
}

// Check what happen if the pass word is not correct
func TestGetUserNotCorrect(t *testing.T) {
	// creating logger
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userRespositoryMock.On(
		"Create", mock.Anything,
	).Return(
		nil,
	)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	user, err := userService.NewUser("bob", "bob@email.com", "1234")

	require.NoError(t, err)
	userRespositoryMock.On(
		"Find", mock.Anything, nil, nil,
	).Return(
		pagination.NewPage([]*models.User{user}, 1, 10, 50),
		nil,
	)

	userFound, err := userService.GetUser(dto.UserLoginDTO{Email: "bob@email.com", Password: "<sdùfjidsfnjd"})
	require.Error(t, err)
	assert.Nil(t, userFound)
}

// Check what happen if the repository return an empty list
func TestGetUserEmpty(t *testing.T) {
	// Creating logger
	observedZapCore, _ := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User, uuid.UUID](t)
	userRespositoryMock.On(
		"Create", mock.Anything,
	).Return(
		nil,
	)
	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	_, err := userService.NewUser("bob", "bob@email.com", "1234")

	require.NoError(t, err)
	userRespositoryMock.On(
		"Find", mock.Anything, nil, nil,
	).Return(
		pagination.NewPage([]*models.User{}, 1, 10, 50),
		nil,
	)

	userFound, err := userService.GetUser(dto.UserLoginDTO{Email: "bob@email.com", Password: "<sdùfjidsfnjd"})
	require.Error(t, err)
	assert.Nil(t, userFound)
}
