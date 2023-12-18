package configuration

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

const (
	ProductionLogger  = "prod"
	DevelopmentLogger = "dev"
)

// The config keys regarding the logger settings
const (
	LoggerModeKey                      string = "logger.mode"
	LoggerDisableStacktraceKey         string = "logger.disableStacktrace"
	LoggerSlowQueryThresholdKey        string = "logger.slowQueryThreshold"
	LoggerSlowTransactionThresholdKey  string = "logger.slowTransactionThreshold"
	LoggerIgnoreRecordNotFoundErrorKey string = "logger.ignoreRecordNotFoundError"
	LoggerParameterizedQueriesKey      string = "logger.parameterizedQueries"
	LoggerRequestTemplateKey           string = "logger.request.template"
)

// Hold the configuration values for the logger
type LoggerConfiguration interface {
	Holder
	GetMode() string
	GetLogLevel() logger.LogLevel
	GetDisableStacktrace() bool
	GetSlowQueryThreshold() time.Duration
	GetSlowTransactionThreshold() time.Duration
	GetIgnoreRecordNotFoundError() bool
	GetParameterizedQueries() bool
	GetRequestTemplate() string
}

// Concrete implementation of the LoggerConfiguration interface
type loggerConfigurationImpl struct {
	mode                      string
	disableStacktrace         bool
	slowQueryThreshold        time.Duration
	slowTransactionThreshold  time.Duration
	ignoreRecordNotFoundError bool
	parameterizedQueries      bool
	requestTemplate           string
}

// Instantiate a new configuration holder for the logger
func NewLoggerConfiguration() LoggerConfiguration {
	loggerConfiguration := new(loggerConfigurationImpl)
	loggerConfiguration.Reload()

	return loggerConfiguration
}

func (loggerConfiguration *loggerConfigurationImpl) Reload() {
	loggerConfiguration.mode = viper.GetString(LoggerModeKey)
	loggerConfiguration.disableStacktrace = viper.GetBool(LoggerDisableStacktraceKey)
	loggerConfiguration.slowQueryThreshold = time.Duration(viper.GetInt(LoggerSlowQueryThresholdKey)) * time.Millisecond
	loggerConfiguration.slowTransactionThreshold = time.Duration(viper.GetInt(LoggerSlowTransactionThresholdKey)) * time.Millisecond
	loggerConfiguration.ignoreRecordNotFoundError = viper.GetBool(LoggerIgnoreRecordNotFoundErrorKey)
	loggerConfiguration.parameterizedQueries = viper.GetBool(LoggerParameterizedQueriesKey)
	loggerConfiguration.requestTemplate = viper.GetString(LoggerRequestTemplateKey)
}

// Return the mode of the logger
func (loggerConfiguration *loggerConfigurationImpl) GetMode() string {
	return loggerConfiguration.mode
}

func (loggerConfiguration *loggerConfigurationImpl) GetLogLevel() logger.LogLevel {
	switch loggerConfiguration.mode {
	case ProductionLogger:
		return logger.Warn
	case DevelopmentLogger:
		return logger.Info
	default:
		return logger.Warn
	}
}

func (loggerConfiguration *loggerConfigurationImpl) GetDisableStacktrace() bool {
	return loggerConfiguration.disableStacktrace
}

func (loggerConfiguration *loggerConfigurationImpl) GetSlowQueryThreshold() time.Duration {
	return loggerConfiguration.slowQueryThreshold
}

func (loggerConfiguration *loggerConfigurationImpl) GetSlowTransactionThreshold() time.Duration {
	return loggerConfiguration.slowTransactionThreshold
}

func (loggerConfiguration *loggerConfigurationImpl) GetIgnoreRecordNotFoundError() bool {
	return loggerConfiguration.ignoreRecordNotFoundError
}

func (loggerConfiguration *loggerConfigurationImpl) GetParameterizedQueries() bool {
	return loggerConfiguration.parameterizedQueries
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
