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

var HTTPServerConfigurationString = `server:
  port: 8000
  host: "0.0.0.0" # listening on all interfaces
  timeout: 15 # in seconds
`

func TestHTTPServerConfigurationNewHttpServerConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewHTTPServerConfiguration(), "the constructor for HttpServerConfiguration should not return a nil value")
}

func TestHTTPServerConfigurationGetPort(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)

	httpServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, 8000, httpServerConfiguration.GetPort())
}

func TestHTTPServerConfigurationGetHost(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)

	httpServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, "0.0.0.0", httpServerConfiguration.GetHost())
}

func TestHTTPServerConfigurationGetAddr(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)

	httpServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, "0.0.0.0:8000", httpServerConfiguration.GetAddr())
}

func TestHTTPServerConfigurationGetMaxTimeout(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)

	httpServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, 15*time.Second, httpServerConfiguration.GetMaxTimeout())
}

func TestHTTPServerConfigurationLog(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	httpServerConfiguration := configuration.NewHTTPServerConfiguration()
	httpServerConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "HTTP Server configuration", log.Message)
	require.Len(t, log.Context, 3)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "port", Type: zapcore.Int64Type, Integer: 8000},
		{Key: "host", Type: zapcore.StringType, String: "0.0.0.0"},
		{Key: "timeout", Type: zapcore.DurationType, Integer: int64(time.Second * 15)},
	}, log.Context)
}
