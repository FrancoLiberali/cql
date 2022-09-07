package configuration

import (
	"time"

	"github.com/spf13/viper"
)

// Hold the configuration values for the http server
type HTTPServerConfiguration struct{}

// Instantiate a new configuration holder for the http server
func NewHTTPServerConfiguration() *HTTPServerConfiguration {
	return &HTTPServerConfiguration{}
}

// Return the host addr
func (hsc *HTTPServerConfiguration) GetHost() string {
	return viper.GetString("host")
}

// Return the port number
func (hsc *HTTPServerConfiguration) GetPort() int {
	return viper.GetInt("port")
}

// Return the maximum timout for read and write
func (hsc *HTTPServerConfiguration) GetMaxTimout() time.Duration {
	return intToSecond(viper.GetInt("max_timeout"))
}
