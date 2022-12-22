package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

func initInitialisationCommands(cfg *verdeter.VerdeterCommand) {

	cfg.GKey(configuration.InitializationDefaultAdminPasswordKey, verdeter.IsStr, "",
		"Set the default admin password is the admin user is not created yet.")
	cfg.SetDefault(configuration.InitializationDefaultAdminPasswordKey, "admin")
}
