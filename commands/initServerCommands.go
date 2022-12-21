package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initServerCommands(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.ServerTimeoutKey, verdeter.IsInt, "", "Maximum timeout of the http server in second (default is 15s)")
	cfg.SetDefault(configuration.ServerTimeoutKey, 15)

	cfg.GKey(configuration.ServerHostKey, verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	cfg.SetDefault(configuration.ServerHostKey, "0.0.0.0")

	cfg.GKey(configuration.ServerPortKey, verdeter.IsInt, "p", "Port to bind (default is 8000)")
	cfg.AddValidator(configuration.ServerPortKey, validators.CheckTCPHighPort)
	cfg.SetDefault(configuration.ServerPortKey, 8000)

	cfg.GKey(configuration.ServerPaginationMaxElemPerPage, verdeter.IsUint, "", "The max number of records returned per page")
	cfg.SetDefault(configuration.ServerPaginationMaxElemPerPage, 100)

}
