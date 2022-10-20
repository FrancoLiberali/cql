package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Hold the configuration values for the logger
type LoggerConfiguration interface {
	ConfigurationHolder
	GetMode() string
	GetRequestTemplate() string
}

// Concrete implementation of the LoggerConfiguration interface
type loggerConfigurationImpl struct {
	mode, requestTemplate string
}

// Instantiate a new configuration holder for the logger
func NewLoggerConfiguration() LoggerConfiguration {
	loggerConfiguration := new(loggerConfigurationImpl)
	loggerConfiguration.Reload()
	return loggerConfiguration
}

func (loggerConfiguration *loggerConfigurationImpl) Reload() {
	loggerConfiguration.mode = viper.GetString("logger.mode")
	loggerConfiguration.requestTemplate = viper.GetString("logger.request.template")
}

// Return the mode of the logger
func (loggerConfiguration *loggerConfigurationImpl) GetMode() string {
	return loggerConfiguration.mode
}

// Return the template string for logging request
func (loggerConfiguration *loggerConfigurationImpl) GetRequestTemplate() string {
	return loggerConfiguration.requestTemplate
}

func (loggerConfiguration *loggerConfigurationImpl) Log() {
	zap.L().Info("Database configuration",
		zap.String("mode", loggerConfiguration.mode),
		zap.String("requestTemplate", loggerConfiguration.requestTemplate),
	)
}
