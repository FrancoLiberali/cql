package commands

// This file holds functions needed by the badaas rootCommand, thoses functions help in creating the http.Server.

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
)

// Create the server from the configuration holder and the http handler
func createServerFromConfigurationHolder(router http.Handler, httpServerConfig configuration.HTTPServerConfiguration) *http.Server {
	address := addrFromConf(httpServerConfig.GetHost(), httpServerConfig.GetPort())
	timeout := httpServerConfig.GetMaxTimeout()
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

func NewHTTPServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	router http.Handler,
	httpServerConfig configuration.HTTPServerConfiguration,
) *http.Server {
	srv := createServerFromConfigurationHolder(router, httpServerConfig)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			logger.Sugar().Infof("Ready to serve at %s", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
