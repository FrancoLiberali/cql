package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// The config keys regarding the logger settings
const (
	LoggerModeKey            string = "logger.mode"
	LoggerRequestTemplateKey string = "logger.request.template"
)

// Hold the configuration values for the logger
type LoggerConfiguration interface {
	Holder
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
	loggerConfiguration.mode = viper.GetString(LoggerModeKey)
	loggerConfiguration.requestTemplate = viper.GetString(LoggerRequestTemplateKey)
}

// Return the mode of the logger
func (loggerConfiguration *loggerConfigurationImpl) GetMode() string {
	return loggerConfiguration.mode
}

// Return the template string for logging request
func (loggerConfiguration *loggerConfigurationImpl) GetRequestTemplate() string {
	return loggerConfiguration.requestTemplate
}

// Log the values provided by the configuration holder
func (loggerConfiguration *loggerConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("Logger configuration",
		zap.String("mode", loggerConfiguration.mode),
		zap.String("requestTemplate", loggerConfiguration.requestTemplate),
	)
}
