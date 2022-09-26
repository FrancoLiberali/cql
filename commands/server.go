package commands

// This file holds functions needed by the badaas rootCommand, thoses functions help in creating the http.Server.

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ditrit/badaas/configuration"
)

// Create an http server from configuration
func createServerFromConfiguration(router http.Handler) *http.Server {
	httpServerConfig := configuration.NewHTTPServerConfiguration()
	return createServerFromConfigurationHolder(router, httpServerConfig)
}

// Create the server from the configuration holder and the http handler
func createServerFromConfigurationHolder(router http.Handler, httpServerConfig *configuration.HTTPServerConfiguration) *http.Server {
	address := addrFromConf(httpServerConfig.GetHost(), httpServerConfig.GetPort())
	timeout := httpServerConfig.GetMaxTimout()
	return createServer(router, address, timeout, timeout)
}

// Create an http server
func createServer(router http.Handler, address string, writeTimeout, readTimeout time.Duration) *http.Server {
	srv := &http.Server{
		Handler: router,
		Addr:    address,

		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	return srv
}

// Create the addr string for the http.Server
// returns "<host>:<port>"
func addrFromConf(host string, port int) string {
	address := fmt.Sprintf("%s:%d",
		host,
		port,
	)
	return address
}
