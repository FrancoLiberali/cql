package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// The config keys regarding the initialization settings
const (
	InitializationDefaultAdminPasswordKey string = "default.admin.password"
)

// Hold the configuration values for the initialization
type InitializationConfiguration interface {
	Holder
	GetAdminPassword() string
}

// Concrete implementation of the InitializationConfiguration interface
type initializationConfigurationIml struct {
	adminPassword string
}

// InitializationConfiguration constructor
func NewInitializationConfiguration() InitializationConfiguration {
	initializationConfiguration := &initializationConfigurationIml{}
	initializationConfiguration.Reload()

	return initializationConfiguration
}

// Reload the InitializationConfiguration
func (initializationConfiguration *initializationConfigurationIml) Reload() {
	initializationConfiguration.adminPassword = viper.GetString(InitializationDefaultAdminPasswordKey)
}

// Log the values provided by the configuration holder
func (initializationConfiguration *initializationConfigurationIml) Log(logger *zap.Logger) {
	logger.Info("Initialization configuration",
		zap.String("adminPassword", initializationConfiguration.adminPassword),
	)
}

// Return default admin password
func (initializationConfiguration *initializationConfigurationIml) GetAdminPassword() string {
	return initializationConfiguration.adminPassword
}
