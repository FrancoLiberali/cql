package commands

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ditrit/badaas/router"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
	"github.com/spf13/viper"
)

var rootCfg = verdeter.NewVerdeterCommand(
	"badaas",

	"Backend and Distribution as a Service",

	`Badaas stands for Backend and Distribution as a Service.`,

	func(cfg *verdeter.VerdeterCommand, args []string) error {
		router := router.SetupRouter()

		address := fmt.Sprintf("%s:%d", viper.Get("host"), viper.GetInt("port"))
		srv := &http.Server{
			Handler: router,
			Addr:    address,

			WriteTimeout: time.Duration(viper.GetInt("max_timeout")) * time.Second,
			ReadTimeout:  time.Duration(viper.GetInt("max_timeout")) * time.Second,
		}

		fmt.Printf("Ready to serve at %s\n", address)
		log.Fatal(srv.ListenAndServe())
		return nil
	},
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
	rootCfg.SetNormalize("config_path", func(val interface{}) interface{} {
		strval, ok := val.(string)
		if ok && strval != "" {
			lastChar := strval[len(strval)-1:]
			if lastChar != "/" {
				return strval + "/"
			}
			return strval
		}

		return nil
	})

	rootCfg.GKey("max_timeout", verdeter.IsInt, "", "maximum timeout (in second)")
	rootCfg.SetDefault("max_timeout", 15)

	rootCfg.GKey("host", verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	rootCfg.SetDefault("host", "0.0.0.0")

	rootCfg.GKey("port", verdeter.IsInt, "p", "Port to bind (default is 8000)")
	rootCfg.SetValidator("port", validators.CheckTCPHighPort)
	rootCfg.SetDefault("port", 8000)
}
