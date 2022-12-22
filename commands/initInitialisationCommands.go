package commands

import (
	"github.com/ditrit/verdeter"
)

func initInitialisationCommands(cfg *verdeter.VerdeterCommand) {

	cfg.GKey("default.admin.password", verdeter.IsStr, "",
		"Set the default admin password is the admin user is not created yet.")
	cfg.SetDefault("default.admin.password", "admin")
}
