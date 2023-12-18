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

var initializationConfigurationString = `default:
  admin:
    password: admin`

func TestInitializationConfigurationInitializationConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewInitializationConfiguration(), "the constructor for InitializationConfiguration should not return a nil value")
}

func TestInitializationConfigurationGetInit(t *testing.T) {
	setupViperEnvironment(initializationConfigurationString)

	initializationConfiguration := configuration.NewInitializationConfiguration()
	assert.Equal(t, "admin", initializationConfiguration.GetAdminPassword())
}

func TestInitializationConfigurationLog(t *testing.T) {
	setupViperEnvironment(initializationConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	initializationConfiguration := configuration.NewInitializationConfiguration()
	initializationConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Initialization configuration", log.Message)
	require.Len(t, log.Context, 1)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "adminPassword", Type: zapcore.StringType, String: "admin"},
	}, log.Context)
}
