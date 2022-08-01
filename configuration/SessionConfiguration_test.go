package configuration_test

import (
	"testing"
	"time"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var SessionConfigurationString = `session:
  duration: 3600 # one hour
  pullInterval: 30 # 30 seconds
  rollDuration: 10 # 10 seconds`

func TestSessionConfigurationNewSessionConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewSessionConfiguration(), "the contructor for PaginationConfiguration should not return a nil value")
}

func TestSessionConfigurationGetSessionDuration(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)
	SessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Duration(time.Hour), SessionConfiguration.GetSessionDuration())
}

func TestSessionConfigurationGetPullIntervall(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)
	SessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Duration(time.Second*30), SessionConfiguration.GetPullInterval())
}

func TestSessionConfigurationGetRollIntervall(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)
	SessionConfiguration := configuration.NewSessionConfiguration()
	assert.Equal(t, time.Duration(time.Second*10), SessionConfiguration.GetRollDuration())
}

func TestSessionConfigurationLog(t *testing.T) {
	setupViperEnvironment(SessionConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	PaginationConfiguration := configuration.NewSessionConfiguration()
	PaginationConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Session configuration", log.Message)
	require.Len(t, log.Context, 3)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "sessionDuration", Type: zapcore.DurationType, Integer: int64(time.Duration(time.Hour))},
		{Key: "pullInterval", Type: zapcore.DurationType, Integer: int64(time.Duration(time.Second * 30))},
		{Key: "rollDuration", Type: zapcore.DurationType, Integer: int64(time.Duration(time.Second * 10))},
	}, log.Context)
}
