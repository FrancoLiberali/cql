package badaas

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
	"github.com/ditrit/verdeter"
)

var BaDaaS = Initializer{}

type Initializer struct {
	modules []fx.Option
}

// Allows to select which modules provided by badaas must be added to the application
func (badaas *Initializer) AddModules(modules ...fx.Option) *Initializer {
	badaas.modules = append(badaas.modules, modules...)

	return badaas
}

// Allows to provide constructors to the application
// so that the constructed objects will be available via dependency injection
func (badaas *Initializer) Provide(constructors ...any) *Initializer {
	badaas.modules = append(badaas.modules, fx.Provide(constructors...))

	return badaas
}

// Allows to invoke functions when the application starts.
// They can take advantage of dependency injection
func (badaas *Initializer) Invoke(funcs ...any) *Initializer {
	badaas.modules = append(badaas.modules, fx.Invoke(funcs...))

	return badaas
}

// Start the application
func (badaas Initializer) Start() {
	rootCommand := verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
		Use:   "badaas",
		Short: "BaDaaS",
		Run:   badaas.runHTTPServer,
	})

	err := configuration.NewCommandInitializer(configuration.NewKeySetter()).Init(rootCommand)
	if err != nil {
		panic(err)
	}

	rootCommand.Execute()
}

// Run the http server for badaas
func (badaas Initializer) runHTTPServer(_ *cobra.Command, _ []string) {
	modules := []fx.Option{
		// internal modules
		configuration.ConfigurationModule,
		router.RouterModule,
		logger.LoggerModule,
		persistence.PersistanceModule,

		// logger for fx
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),

		// create httpServer
		fx.Provide(newHTTPServer),
		// Finally: we invoke the newly created server
		fx.Invoke(func(*http.Server) { /* we need this function to be empty*/ }),
	}

	fx.New(
		// add modules selected by user
		append(modules, badaas.modules...)...,
	).Run()
}
