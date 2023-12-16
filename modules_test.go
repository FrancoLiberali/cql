package badaas

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	mocksConfiguration "github.com/ditrit/badaas/mocks/configuration"
	mockUserServices "github.com/ditrit/badaas/mocks/services/userservice"
)

func TestCreateSuperUser(t *testing.T) {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocksConfiguration.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, nil)
	err := createSuperUser(
		initializationConfig,
		logger,
		userService,
	)
	assert.NoError(t, err)
}

func TestCreateSuperUser_UserExists(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocksConfiguration.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, errors.New("user already exist in database"))
	err := createSuperUser(
		initializationConfig,
		logger,
		userService,
	)
	assert.NoError(t, err)

	require.Equal(t, 1, logs.Len())
}

func TestCreateSuperUser_UserServiceError(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	initializationConfig := mocksConfiguration.NewInitializationConfiguration(t)
	initializationConfig.On("GetAdminPassword").Return("adminpassword")
	userService := mockUserServices.NewUserService(t)
	userService.
		On("NewUser", "admin", "admin-no-reply@badaas.com", "adminpassword").
		Return(nil, errors.New("email not valid"))
	err := createSuperUser(
		initializationConfig,
		logger,
		userService,
	)
	assert.Error(t, err)

	require.Equal(t, 1, logs.Len())
}
