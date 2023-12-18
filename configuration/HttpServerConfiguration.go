package configuration

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/utils"
)

// The config keys regarding the http server settings
const (
	ServerTimeoutKey               string = "server.timeout"
	ServerHostKey                  string = "server.host"
	ServerPortKey                  string = "server.port"
	ServerPaginationMaxElemPerPage string = "server.pagination.page.max"
)

// Hold the configuration values for the http server
type HTTPServerConfiguration interface {
	Holder
	GetAddr() string
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
	httpServerConfiguration.host = viper.GetString(ServerHostKey)
	httpServerConfiguration.port = viper.GetInt(ServerPortKey)
	httpServerConfiguration.timeout = utils.IntToSecond(viper.GetInt(ServerTimeoutKey))
}

// Return the host addr
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetHost() string {
	return httpServerConfiguration.host
}

// Return the port number
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetPort() int {
	return httpServerConfiguration.port
}

// Return the maximum timeout for read and write
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetMaxTimeout() time.Duration {
	return httpServerConfiguration.timeout
}

// Log the values provided by the configuration holder
func (httpServerConfiguration *hTTPServerConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("HTTP Server configuration",
		zap.String("host", httpServerConfiguration.host),
		zap.Int("port", httpServerConfiguration.port),
		zap.Duration("timeout", httpServerConfiguration.timeout),
	)
}

// Create the addr string in format: "<host>:<port>"
func (httpServerConfiguration *hTTPServerConfigurationImpl) GetAddr() string {
	return fmt.Sprintf("%s:%d",
		httpServerConfiguration.GetHost(),
		httpServerConfiguration.GetPort(),
	)
}
