package configuration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ditrit/badaas/configuration"
)

var LoggerConfigurationString = `logger:
  mode: prod
  request:
    template: "{proto} {method} {url}"
`

func TestLoggerConfigurationNewLoggerConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewLoggerConfiguration(), "the constructor for LoggerConfiguration should not return a nil value")
}

func TestLoggerConfigurationLoggerGetMode(t *testing.T) {
	setupViperEnvironment(LoggerConfigurationString)

	loggerConfiguration := configuration.NewLoggerConfiguration()
	assert.Equal(t, "prod", loggerConfiguration.GetMode())
}

func TestLoggerConfigurationLoggerRequestTemplate(t *testing.T) {
	setupViperEnvironment(LoggerConfigurationString)

	loggerConfiguration := configuration.NewLoggerConfiguration()
	assert.Equal(t, "{proto} {method} {url}", loggerConfiguration.GetRequestTemplate())
}

func TestLoggerConfigurationLog(t *testing.T) {
	setupViperEnvironment(LoggerConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	loggerConfiguration := configuration.NewLoggerConfiguration()
	loggerConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Logger configuration", log.Message)
	require.Len(t, log.Context, 2)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "requestTemplate", Type: zapcore.StringType, String: "{proto} {method} {url}"},
		{Key: "mode", Type: zapcore.StringType, String: "prod"},
	}, log.Context)
}
