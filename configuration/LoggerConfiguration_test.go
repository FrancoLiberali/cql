package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
)

var LoggerConfigurationString = `logger:
  mode: prod
  request:
    template: "{proto} {method} {url}"
`

func TestLoggerConfigurationNewLoggerConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewLoggerConfiguration(), "the contructor for LoggerConfiguration should not return a nil value")
}

func TestLoggerConfigurationLoggerGetMode(t *testing.T) {
	setupViperEnvironment(LoggerConfigurationString)
	LoggerConfiguration := configuration.NewLoggerConfiguration()
	assert.Equal(t, "prod", LoggerConfiguration.GetMode())
}

func TestLoggerConfigurationLoggerRequestTemplate(t *testing.T) {
	setupViperEnvironment(LoggerConfigurationString)
	LoggerConfiguration := configuration.NewLoggerConfiguration()
	assert.Equal(t, "{proto} {method} {url}", LoggerConfiguration.GetRequestTemplate())
}