package commands

import (
	"log"

	"github.com/ditrit/badaas/logger"
	"github.com/ditrit/badaas/persistence/registry"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/services/session"
	"github.com/ditrit/badaas/services/userservice"
	"github.com/ditrit/verdeter"
	"go.uber.org/zap"
)

// Create a super admin user and exit with code 1 on error
func createSuperAdminUser() {
	logg := zap.L().Sugar()
	_, err := userservice.NewUser("superadmin", "superadmin@badaas.test", "1234")
	if err != nil {
		if repository.ErrAlreadyExists == err {
			logg.Debugf("The superadmin user already exists in database")
		} else {
			logg.Fatalf("failed to save the super admin %w", err)
		}
	}

}

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	err := logger.InitLoggerFromConf()
	if err != nil {
		log.Fatalf("An error happened while initializing logger (ERROR=%s)", err.Error())
	}

	zap.L().Info("The logger is initialiazed")

	// create router
	router := router.SetupRouter()

	registryInstance, err := registry.FactoryRegistry(registry.GormDataStore)
	if err != nil {
		zap.L().Sugar().Fatalf("An error happened while initializing datastorage layer (ERROR=%s)", err.Error())
	}
	registry.ReplaceGlobals(registryInstance)
	zap.L().Info("The datastorage layer is initialized")

	createSuperAdminUser()

	err = session.Init()
	if err != nil {
		zap.L().Sugar().Fatalf("An error happened while initializing the session service (ERROR=%s)", err.Error())
	}
	zap.L().Info("The session service is initialized")

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
	rootCfg.Initialize()

	rootCfg.GKey("config_path", verdeter.IsStr, "", "Path to the config file/directory")
	rootCfg.SetDefault("config_path", ".")

	initServerCommands(rootCfg)
	initLoggerCommands(rootCfg)
	initDatabaseCommands(rootCfg)
}
