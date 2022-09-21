package configuration

import "github.com/spf13/viper"

// Hold the configuration values for the logger
type LoggerConfiguration struct{}

// Instantiate a new configuration holder for the logger
func NewLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{}
}

// Return the mode of the logger
func (lc *LoggerConfiguration) GetMode() string {
	return viper.GetString("logger.mode")
}
