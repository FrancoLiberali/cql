package commands

import (
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

// Run the http server for badaas
func runHTTPServer(cfg *verdeter.VerdeterCommand, args []string) error {
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
}
