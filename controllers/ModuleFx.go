package controllers

import (
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/services/userservice"
)

var InfoControllerModule = fx.Module(
	"infoController",
	fx.Provide(NewInfoController),
	fx.Invoke(AddInfoRoutes),
)

var AuthControllerModule = fx.Module(
	"authController",
	fx.Provide(NewBasicAuthenticationController),
	fx.Invoke(AddAuthRoutes),
	fx.Invoke(createSuperUser),
)

// Create a super user
func createSuperUser(
	config configuration.InitializationConfiguration,
	logger *zap.Logger,
	userService userservice.UserService,
) error {
	// Create a super admin user and exit with code 1 on error
	_, err := userService.NewUser("admin", "admin-no-reply@badaas.com", config.GetAdminPassword())
	if err != nil {
		if !strings.Contains(err.Error(), "already exist in database") {
			logger.Sugar().Errorf("failed to save the super admin %w", err)
			return err
		}
		logger.Sugar().Infof("The superadmin user already exists in database")
	}
	return nil
}
