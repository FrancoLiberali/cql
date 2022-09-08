package commands

import (
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/db"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
	err := db.InitializeDBFromConf()
	if err != nil {
		return fmt.Errorf("failed to initialize the connection to the database, (%w)", err)
	}
	err = db.AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate the database, (%w)", err)
	}

	// configuration holder for the http server
	// get the config value with the correct types by using the method on this struct
	httpServerConfig := configuration.NewHTTPServerConfiguration()

	router := router.SetupRouter()

	address := fmt.Sprintf("%s:%d",
		httpServerConfig.GetHost(),
		httpServerConfig.GetPort(),
	)
	srv := &http.Server{
		Handler: router,
		Addr:    address,

		WriteTimeout: httpServerConfig.GetMaxTimout(),
		ReadTimeout:  httpServerConfig.GetMaxTimout(),
	}

	fmt.Printf("Ready to serve at %s\n", address)
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

	rootCfg.GKey("max_timeout", verdeter.IsInt, "", "maximum timeout (in second)")
	rootCfg.SetDefault("max_timeout", 15)

	rootCfg.GKey("host", verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	rootCfg.SetDefault("host", "0.0.0.0")

	rootCfg.GKey("port", verdeter.IsInt, "p", "Port to bind (default is 8000)")
	rootCfg.SetValidator("port", validators.CheckTCPHighPort)
	rootCfg.SetDefault("port", 8000)

	rootCfg.GKey("database.port", verdeter.IsInt, "", "[DB] the port of the database server")
	rootCfg.SetRequired("database.port")

	rootCfg.GKey("database.host", verdeter.IsStr, "", "[DB] the host of the database server")
	rootCfg.SetRequired("database.host")

	rootCfg.GKey("database.name", verdeter.IsStr, "", "[DB] the name of the database to use")
	rootCfg.SetRequired("database.name")

	rootCfg.GKey("database.username", verdeter.IsStr, "", "[DB] the username of the account on the database server")
	rootCfg.SetRequired("database.username")

	rootCfg.GKey("database.password", verdeter.IsStr, "", "[DB] the password of the account one the database server")
	rootCfg.SetRequired("database.password")

	rootCfg.GKey("database.sslmode", verdeter.IsStr, "", "[DB] the sslmode to use when connecting to the database server")
	rootCfg.SetRequired("database.sslmode")

}
