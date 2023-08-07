package configuration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ditrit/badaas/configuration"
)

var SessionConfigurationString = `session:
  duration: 3600 # one hour
  pullInterval: 30 # 30 seconds
  rollDuration: 10 # 10 seconds`

func TestSessionConfigurationNewSessionConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewSessionConfiguration(), "the constructor for PaginationConfiguration should not return a nil value")
}

func TestSessionConfigurationGetSessionDuration(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)

	sessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Hour, sessionConfiguration.GetSessionDuration())
}

func TestSessionConfigurationGetPullInterval(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)

	sessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Second*30, sessionConfiguration.GetPullInterval())
}

func TestSessionConfigurationGetRollInterval(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)

	sessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Second*10, sessionConfiguration.GetRollDuration())
}

func TestSessionConfigurationLog(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	paginationConfiguration := configuration.NewSessionConfiguration()
	paginationConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Session configuration", log.Message)
	require.Len(t, log.Context, 3)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "sessionDuration", Type: zapcore.DurationType, Integer: int64(time.Hour)},
		{Key: "pullInterval", Type: zapcore.DurationType, Integer: int64(time.Second * 30)},
		{Key: "rollDuration", Type: zapcore.DurationType, Integer: int64(time.Second * 10)},
	}, log.Context)
}
