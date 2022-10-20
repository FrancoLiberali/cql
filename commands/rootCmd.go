package commands

import (
	"log"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence/registry"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"go.uber.org/zap"
)

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	conf := configuration.NewConfiguration()
	configuration.ReplaceGlobals(conf)
	err := logger.InitLoggerFromConf()
	if err != nil {
		log.Fatalf("An error happened while initializing logger (ERROR=%s)", err.Error())
	}
	configuration.Get().Log()

	zap.L().Info("The logger is initialiazed")

	// create router
	router := router.SetupRouter()

	registryInstance, err := registry.FactoryRegistry(registry.GormDataStore)
	if err != nil {
		zap.L().Sugar().Fatalf("An error happened while initializing datastorage layer (ERROR=%s)", err.Error())
	}
	registry.ReplaceGlobals(registryInstance)
	zap.L().Info("The datastorage layer is initialized")

	// create server
	srv := createServerFromConfiguration(router)

	zap.L().Sugar().Infof("Ready to serve at %s\n", srv.Addr)
	return srv.ListenAndServe()
}

var rootCfg = verdeter.NewVerdeterCommand(
	"badaas",
	"Backend and Distribution as a Service",
	`Badaas stands for Backend and Distribution as a Service.`,
	runHTTPServer,
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCfg.Execute()
}

func init() {
	rootCfg.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	rootCfg.SetDefault("config_path", ".")

	initServerCommands(rootCfg)
	initLoggerCommands(rootCfg)
	initDatabaseCommands(rootCfg)
}
