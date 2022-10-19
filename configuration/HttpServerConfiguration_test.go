package configuration_test

import (
	"testing"
	"time"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
)

var HTTPServerConfigurationString = `server:
  port: 8000
  host: "0.0.0.0" # listening on all interfaces
  max_timeout: 15 # in seconds
`

func TestHTTPServerConfigurationNewHttpServerConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewHTTPServerConfiguration(), "the contructor for HttpServerConfiguration should not return a nil value")
}

func TestHTTPServerConfigurationGetPort(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)
	HTTPServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, 8000, HTTPServerConfiguration.GetPort())
}
func TestHTTPServerConfigurationGetHost(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)
	HTTPServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, "0.0.0.0", HTTPServerConfiguration.GetHost())
}

func TestHTTPServerConfigurationGetMaxTimeout(t *testing.T) {
	setupViperEnvironment(HTTPServerConfigurationString)
	HTTPServerConfiguration := configuration.NewHTTPServerConfiguration()
	assert.Equal(t, time.Duration(15*time.Second), HTTPServerConfiguration.GetMaxTimeout())
}
