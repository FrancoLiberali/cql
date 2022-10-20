package configuration

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Hold the configuration values for the http server
type HTTPServerConfiguration interface {
	ConfigurationHolder
	GetHost() string
	GetPort() int
	GetMaxTimeout() time.Duration
}

// Concrete implementation of the HTTPServerConfiguration interface
type hTTPServerConfigurationImpl struct {
	host    string
	port    int
	timeout time.Duration
}

// Instantiate a new configuration holder for the http server
func NewHTTPServerConfiguration() HTTPServerConfiguration {
	httpServerConfiguration := new(hTTPServerConfigurationImpl)
	httpServerConfiguration.Reload()
	return httpServerConfiguration
}

// Reload HTTP Server configuration
func (httpServerConfiguration *hTTPServerConfigurationImpl) Reload() {
	httpServerConfiguration.host = viper.GetString("server.host")
	httpServerConfiguration.port = viper.GetInt("server.port")
	httpServerConfiguration.timeout = intToSecond(viper.GetInt("server.max_timeout"))
}

// Return the host addr
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetHost() string {
	return httpServerConfiguration.host
}

// Return the port number
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetPort() int {
	return httpServerConfiguration.port
}

// Return the maximum timout for read and write
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetMaxTimeout() time.Duration {
	return httpServerConfiguration.timeout
}

// Return the host addr
func (httpServerConfiguration *hTTPServerConfigurationImpl) Log() {
	zap.L().Info("HTTP Server configuration",
		zap.String("host", httpServerConfiguration.host),
		zap.Int("port", httpServerConfiguration.port),
		zap.Duration("timeout", httpServerConfiguration.timeout),
	)
}
