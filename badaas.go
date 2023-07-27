package main

import (
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/ditrit/verdeter"
)

// Badaas application, run a http-server on 8000.
func main() {
	rootCommand := verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
		Use:   "badaas",
		Short: "BaDaaS",
		Run:   runHTTPServer,
	})

	err := configuration.NewCommandInitializer(configuration.NewKeySetter()).Init(rootCommand)
	if err != nil {
		panic(err)
	}

	rootCommand.Execute()
}

// Run the http server for badaas
func runHTTPServer(cmd *cobra.Command, args []string) {
	fx.New(
		// Modules
		configuration.ConfigurationModule,
		router.RouterModule,
		logger.LoggerModule,
		persistence.PersistanceModule,

		fx.Provide(userservice.NewUserService),
		fx.Provide(sessionservice.NewSessionService),
		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		fx.Provide(newHTTPServer),

		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
	).Run()
}
