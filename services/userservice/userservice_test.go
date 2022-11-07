package userservice_test

import (
	"testing"

	"github.com/ditrit/badaas/httperrors"
	repositorymocks "github.com/ditrit/badaas/mocks/persistence/repository"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/services/userservice"
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

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User](t)
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

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User](t)
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

	userRespositoryMock := repositorymocks.NewCRUDRepository[models.User](t)

	userService := userservice.NewUserService(observedLogger, userRespositoryMock)
	user, err := userService.NewUser("bob", "bob@", "1234")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Checking logs
	assert.Equal(t, 0, observedLogs.Len())
}
