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

type hTTPServerConfigurationImpl struct {
	host    string
	port    int
	timeout time.Duration
}

// Instantiate a new configuration holder for the http server
func NewHTTPServerConfiguration() HTTPServerConfiguration {
	hsc := new(hTTPServerConfigurationImpl)
	hsc.Reload()
	return hsc
}

// Return the host addr
func (hsc *hTTPServerConfigurationImpl) Reload() {
	hsc.host = viper.GetString("server.host")
	hsc.port = viper.GetInt("server.port")
	hsc.timeout = intToSecond(viper.GetInt("server.max_timeout"))
}

// Return the host addr
func (hsc *hTTPServerConfigurationImpl) GetHost() string {
	return hsc.host
}

// Return the port number
func (hsc *hTTPServerConfigurationImpl) GetPort() int {
	return hsc.port
}

// Return the maximum timout for read and write
func (hsc *hTTPServerConfigurationImpl) GetMaxTimeout() time.Duration {
	return hsc.timeout
}

// Return the host addr
func (hsc *hTTPServerConfigurationImpl) Log() {
	zap.L().Info("HTTP Server configuration",
		zap.String("host", hsc.host),
		zap.Int("port", hsc.port),
		zap.Duration("timeout", hsc.timeout),
	)
}
