package commands

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initServerCommands(cfg *verdeter.VerdeterCommand) {
	serverTimeoutKey := "server.timeout"
	cfg.GKey(serverTimeoutKey, verdeter.IsInt, "", "Maximum timeout of the http server in second (default is 15s)")
	cfg.SetDefault(serverTimeoutKey, 15)

	serverHostKey := "server.host"
	cfg.GKey(serverHostKey, verdeter.IsStr, "", "Address to bind (default is 0.0.0.0)")
	cfg.SetDefault(serverHostKey, "0.0.0.0")

	serverPortKey := "server.port"
	cfg.GKey(serverPortKey, verdeter.IsInt, "p", "Port to bind (default is 8000)")
	cfg.SetValidator(serverPortKey, validators.CheckTCPHighPort)
	cfg.SetDefault(serverPortKey, 8000)

	paginationMaxElemPerPage := "server.pagination.page.max"
	cfg.GKey(paginationMaxElemPerPage, verdeter.IsUint, "", "The max number of records returned per page")
	cfg.SetDefault(paginationMaxElemPerPage, 100)

}
